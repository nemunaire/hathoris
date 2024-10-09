import { CustomSource } from '$lib/custom_source.js';

export class Settings {
  constructor(res) {
    if (res) {
      this.update(res);
    }
  }

  update({ custom_sources }) {
    if (custom_sources) {
      this.custom_sources = custom_sources.map((e) => new CustomSource(e));
    } else {
      this.custom_sources = [];
    }
  }

  async save(avoidStructUpdate) {
    const res = await fetch('api/settings', {
      headers: {'Accept': 'application/json'},
      method: 'POST',
      body: JSON.stringify(this),
    });
    if (res.status == 200) {
      const data = await res.json();
      if (!avoidStructUpdate) {
        if (data == null) {
          this.update({});
        } else {
          this.update(data);
        }
        return this;
      } else {
        return new Settings(data);
      }
    } else {
      throw new Error((await res.json()).errmsg);
    }
  }
}

export async function getSettings() {
  const res = await fetch(`api/settings`, {headers: {'Accept': 'application/json'}})
  if (res.status == 200) {
    const data = await res.json();
    if (data == null) {
      return {}
    } else {
      return new Settings(data);
    }
  } else {
    throw new Error((await res.json()).errmsg);
  }
}
