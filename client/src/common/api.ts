import type { Character, QueryOptions, Weapon } from 'genshin-db';

export const queryCharacter = async (query: string, opts?: QueryOptions) => {
  console.debug(opts);
  const res = await fetch(
    `https://genshin-db-api.vercel.app/api/v5/characters?query=${query}&resultLanguage=ChineseSimplified`,
  );
  return res.json() as Promise<Character>;
};

export const queryWeapon = async (query: string, opts?: QueryOptions) => {
  console.debug(opts);
  const res = await fetch(
    `https://genshin-db-api.vercel.app/api/v5/weapons?query=${query}&resultLanguage=ChineseSimplified`,
  );
  return res.json() as Promise<Weapon>;
};
