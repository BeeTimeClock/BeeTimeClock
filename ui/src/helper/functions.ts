export function autoImplement<T extends object>(defaults: Partial<T> = {}) {
  // eslint-disable-next-line @typescript-eslint/no-extraneous-class
  return class {
    constructor(data: Partial<T> = {}) {
      Object.assign(this, defaults);
      Object.assign(this, data);
    }
  } as new (data?: T) => T;
}
