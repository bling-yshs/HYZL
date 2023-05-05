import plugin from "../../lib/plugins/plugin.js";
import cfg from "../../lib/config/config.js";

let coolingTimeOfVotingKick = 600; //投票踢人CD，单位是秒
let numberOfVotingKicks = 3; //投票踢人的条数，大于等于3就自动踢人
let whiteList = "450977050|377178599"; //白名单QQ号，用|分割，防止有人恶意踢出其它机器人

let keyWords = "快手|广告|代肝"; //触发投票提示的关键词，用|分割

export class example extends plugin {
    constructor() {
        super({
            name: "投票踢人",
            dsc: "简单开发示例",
            event: "message",
            priority: 5000,
            rule: [
                {
                    reg: "^投票踢人",
                    fnc: "kickPeople",
                },
                {
                    reg: keyWords,
                    fnc: "remind",
                },
            ],
        });
    }
    async remind(e) {
        if (!e.group.is_owner && !e.group.is_admin) {
            return false;
        }
        let key = `Yunzai:remind:${e.group_id}`;
        let time = await redis.get(key);
        if (time) {
            return false;
        }
        e.reply(
            `检查到疑似广告的关键词，群友可以发送投票踢人@某人来进行踢人，如： 投票踢人@bling丶一闪\n\n友情提醒：如果滥用投票踢人，会被踢出群聊`
        );
        await redis.set(key, "1", { EX: 60 });
        return true;
    }
    async kickPeople(e) {
        //自己不是群主或者管理员没法踢人
        if (!e.group.is_owner && !e.group.is_admin) {
            return false;
        }
        //别踢主人
        if (e.at == cfg.masterQQ) {
            e.reply("居然想踢主人，岂有此理ヽ(｀⌒´メ)ノ");
            return true;
        }

        let msg = e.message.find((item) => item.type == "at");
        if (!msg) {
            e.reply("快把广告狗@出来，我要踢他！！！");
            return false;
        }
        let atQQ = msg.qq;

        if (whiteList.indexOf(atQQ) != -1) {
            e.reply("此用户在白名单中，无法踢出");
            return false;
        }

        //别踢自己
        if (atQQ == e.self_id) {
            e.reply(`我踢我自己是吧`);
            return false;
        }

        let theObjectOfAtQQ = e.group.pickMember(atQQ);

        if (theObjectOfAtQQ.is_owner) {
            e.reply(`群主你都敢踢╭(°A°\`)╮`);
            return false;
        }

        if (theObjectOfAtQQ.is_admin) {
            e.reply(`管理员你都敢踢╭(°A°\`)╮`);
            return false;
        }

        if (e.member.is_owner || e.member.is_admin) {
            e.group.kickMember(atQQ);
            e.reply(`神权发动成功，${atQQ}已被踢出群聊`);
            return true;
        }

        let qqOfKicker = e.user_id; //踢人的人的QQ号

        if (qqOfKicker == e.at) {
            e.reply(`哥们不至于不至于`);
            return false;
        }

        let iskicked = await redis.get(
            `Yunzai:kick:${e.group_id}:${qqOfKicker}`
        );
        if (iskicked) {
            e.reply(`你已经投过票了`);
            return false;
        }
        await redis.set(`Yunzai:kick:${e.group_id}:${qqOfKicker}`, 1, {
            EX: 120,
        });

        let key = `Yunzai:kickPeople:${e.group_id}`; //redis key
        let time = await redis.get(key); //获取redis中key存不存在
        //如果key还存在，说明还在CD
        if (time) {
            let timeRemaining = await redis.ttl(key); //获取redis里的key的过期时间
            e.reply(`投票踢人还有 ${timeRemaining} 秒 CD`);
            return false;
        }

        let kickKey = `Yunzai:kickPeople:${e.group_id}:${atQQ}`; //redis key

        let kickingNum = 0;
        kickingNum = await redis.get(kickKey); //获取redis中key存不存在
        if (kickingNum >= numberOfVotingKicks - 1) {
            e.group.kickMember(atQQ);
            e.reply(`投票踢人成功，${atQQ}已被踢出群聊`);
            await redis.del(kickKey);
            await redis.set(key, 1, { EX: coolingTimeOfVotingKick });
            return true;
        }
        if (kickingNum) {
            await redis.set(kickKey, ++kickingNum, { EX: 120 });
        } else {
            await redis.set(kickKey, 1, { EX: 120 });
            kickingNum = 1;
        }

        //如果投票数小于3，就提示
        e.reply(
            `正在踢出${e.at}，当前票数 ${kickingNum} ，到达 ${numberOfVotingKicks} 票将自动踢出`
        );
        return true;
    }
}
