import classNames from 'classnames';
import type { Character, QueryOptions, Weapon } from 'genshin-db';
import React, { useEffect, useState } from 'react';
import { queryCharacter, queryWeapon } from '../../common/api';
import { toEnElement } from '../../utils';
import RarityIndicator from '../RarityIndicator';
import styles from './index.module.less';

const queryOption = {
  matchAliases: false,
  matchCategories: false,
  verboseCategories: false,
  queryLanguages: ['English'],
  resultLanguage: 'CHS',
};

const isCharacter = (item: Character | Weapon): item is Character => {
  if (Object.keys(item).includes('gender')) {
    return true;
  }
  return false;
};

const dateTimeFormatter = Intl.DateTimeFormat('zh', {
  timeZone: 'Asia/Shanghai',
  year: 'numeric',
  month: 'numeric',
  day: 'numeric',
});

const ItemName: React.FC<{ itemInfo: Character | Weapon | null }> = props => {
  const { itemInfo } = props;
  if (!itemInfo) {
    return <span></span>;
  }
  return (
    <div
      className={classNames(
        styles.name,
        isCharacter(itemInfo)
          ? styles[`${toEnElement(itemInfo.elementText)}Text`]
          : '',
      )}
    >{`${itemInfo.name}(${
      isCharacter(itemInfo) ? itemInfo.elementText : itemInfo.weaponType
    })`}</div>
  );
};

interface AvatarProps {
  itemInfo: Character | Weapon;
  rarity: string;
}
export const Avatar = (props: AvatarProps) => {
  return (
    <div className={styles.avatar}>
      <span
        className={classNames(
          styles.iconWrapper,
          styles[`star${props.rarity}Bg`],
        )}
      >
        <img
          className={styles.icon}
          src={props.itemInfo.images.mihoyo_icon}
          alt="item-icon"
        />
      </span>
      <RarityIndicator rarity={parseInt(props.rarity)} />
    </div>
  );
};

interface ItemCardProps {
  itemId: string;
  pityStar4: number;
  pityStar5: number;
  time: string;
  itemType: 'weapon' | 'character';
  rarity: string;
}

const ItemCard: React.FC<ItemCardProps> = props => {
  const [itemInfo, setItemInfo] = useState<Character | Weapon>();

  useEffect(() => {
    const getInfo = async () => {
      if (props.itemType === 'character') {
        const itemInfo = await queryCharacter(
          props.itemId,
          queryOption as unknown as QueryOptions,
        );
        setItemInfo(itemInfo);
        return;
      } else {
        const itemInfo = await queryWeapon(
          props.itemId,
          queryOption as unknown as QueryOptions,
        );
        setItemInfo(itemInfo);
        return;
      }
    };
    getInfo();
  }, [props.itemId, props.itemType]);

  if (!itemInfo) {
    return <div />;
  }

  return (
    <li className={styles.item}>
      <Avatar itemInfo={itemInfo} rarity={props.rarity} />
      <ItemName itemInfo={itemInfo} />
      <div className={styles.datetime}>
        {dateTimeFormatter.format(new Date(props.time))}
      </div>
      {props.rarity === '4' ? (
        <div className={styles.pity}>{props.pityStar4}</div>
      ) : null}
      {props.rarity === '5' ? (
        <div className={styles.pity}>{props.pityStar5}</div>
      ) : null}
    </li>
  );
};

export default ItemCard;
