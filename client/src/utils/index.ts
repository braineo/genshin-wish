const elements: { [id: string]: string } = {
  冰: 'cryo',
  火: 'pyro',
  岩: 'geo',
  草: 'dendro',
  水: 'hydro',
  风: 'anemo',
  雷: 'electro',
};

export const toEnElement = (cnElement: string): string | undefined => {
  return elements[cnElement];
};
