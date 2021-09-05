/// <reference types="react-scripts" />

declare module '*.module.less' {
  const classes: { [key: string]: string };
  export default classes;
}

declare module 'genshin-wish' {
  export type WishLog = {
    gachaType: string;
    id: string;
    Item: GachaItem;
    time: string;
    pityStar4: number;
    pityStar5: number;
  };
  export type GachaItem = {
    id: string;
    name: string;
    type: 'weapon' | 'character';
    rarity: string;
  };
}
