import React from 'react';
import starSvg from '../../assets/star.svg';
import styles from './index.module.less';

type RarityIndicatorProps = {
  rarity: number;
};
const RarityIndicator: React.FC<RarityIndicatorProps> = props => {
  return (
    <div className={styles.rarity}>
      {Array.from(Array(props.rarity)).map((_, index) => (
        <img key={index} className={styles.star} src={starSvg} alt="star" />
      ))}
    </div>
  );
};

export default RarityIndicator;
