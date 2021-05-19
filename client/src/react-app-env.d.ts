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

declare module 'genshin-db' {
  export type ConfigOptions = {
    matchAliases: boolean; // Allows the matching of aliases.
    matchCategories: boolean; // Allows the matching of categories. If true, then returns an array if it matches.
    verboseCategories: boolean; // Used if a category is matched. If true, then replaces each string name in the array with the data object instead.
    queryLanguages: string[]; // Array of languages that your query will be searched in.
    resultLanguage: string; // Output language that you want your results to be in.
  };

  export type Character = {
    name: string;
    title: string;
    description: string;
    rarity: string;
    element: string;
    weapontype: string;
    substat: string;
    gender: string;
    body: string;
    association: string;
    region: string;
    affiliation: string;
    birthdaymmdd: string;
    birthday: string;
    constellation: string;
    cv: {
      english: string;
      chinese: string;
      japanese: string;
      korean: string;
    };
    talentmaterialtype: string;
    images: {
      image?: string;
      card?: string;
      portrait?: string;
      icon: string;
      sideicon: string;
      cover1?: string;
      cover2?: string;
      'hoyolab-avatar'?: string;
    };
    url: { fandom: string };
  };

  export type Weapon = {
    name: string;
    description: string;
    weapontype: string;
    rarity: string;
    baseatk: 48;
    substat: string;
    subvalue: string;
    effectname: string;
    effect: string;
    r1: string[];
    r2: string[];
    r3: string[];
    r4: string[];
    r5: string[];
    weaponmaterialtype: string;
    images: {
      image: string;
      icon: string;
      awakenicon: string;
    };
    url: { fandom: string };
  };

  export function setOptions(opts: ConfigOptions): void;
  export function character(query: string, opts?: ConfigOptions): ?Character;
  export function weapon(query: string, opts?: ConfigOptions): ?Weapon;
}
