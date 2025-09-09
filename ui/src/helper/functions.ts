export function autoImplement<T extends object>(defaults: Partial<T> = {}) {
  return class {
    constructor(data: Partial<T> = {}) {
      Object.assign(this, defaults);
      Object.assign(this, data);
    }
  } as new (data?: T) => T;
}
