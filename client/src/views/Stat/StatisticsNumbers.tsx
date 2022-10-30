import { Card, Row, Statistic } from 'antd';
import { WishLog } from 'genshin-wish';

type WishStatisticsProps = {
  wishLogs: WishLog[];
};
const StatisticsNumbers: React.FC<WishStatisticsProps> = props => {
  const { wishLogs } = props;
  const star5Wishes = wishLogs.filter(log => log.Item.rarity === '5');

  return (
    <Row>
      <Card>
        <Statistic
          title="距上次五星"
          value={
            wishLogs.length > 0 && wishLogs[0].Item.rarity !== '5'
              ? wishLogs[0].pityStar5
              : 0
          }
        />
      </Card>
      <Card>
        <Statistic
          title="距上次四星"
          value={
            wishLogs.length > 0 && wishLogs[0].Item.rarity !== '4'
              ? wishLogs[0].pityStar4
              : 0
          }
        />
      </Card>
      <Card>
        <Statistic title="祈愿数" value={wishLogs.length} />
      </Card>
      <Card>
        <Statistic title="五星物品总数" value={star5Wishes.length} />
      </Card>
      <Card>
        <Statistic
          title="五星平均间隔"
          precision={1}
          value={
            star5Wishes
              .map(log => log.pityStar5)
              .reduce((acc, current) => acc + current, 0) / star5Wishes.length
          }
        />
      </Card>
      <Card>
        <Statistic
          title="四星物品总数"
          value={wishLogs.filter(log => log.Item.rarity === '4').length}
        />
      </Card>
    </Row>
  );
};

export default StatisticsNumbers;
