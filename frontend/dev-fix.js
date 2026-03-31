const { Event } = await import('node:events');

globalThis.CustomEvent = class CustomEvent extends Event {
  constructor(type, options) {
    super(type, options);
    if (options && options.detail !== undefined) {
      this.detail = options.detail;
    }
  }
};

await import('./node_modules/vite/bin/vite.js');
