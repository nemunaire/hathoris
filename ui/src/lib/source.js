export class Source {
  constructor(id, res) {
    this.id = id;
    if (res) {
      this.update(res);
    }
  }

  update({ name, enabled, active }) {
    this.name = name;
    this.enabled = enabled;
    this.active = active;
  }

  async activate() {
    await fetch(`api/sources/${this.id}/enable`, {headers: {'Accept': 'application/json'}, method: 'POST'});
  }

  async deactivate() {
    await fetch(`api/sources/${this.id}/disable`, {headers: {'Accept': 'application/json'}, method: 'POST'});
  }

  async currently() {
    const data = await fetch(`api/sources/${this.id}/currently`, {headers: {'Accept': 'application/json'}});
    if (data.status == 200) {
      return await data.json();
    } else {
      throw new Error((await res.json()).errmsg);
    }
  }
}

export async function getSources() {
  const res = await fetch(`api/sources`, {headers: {'Accept': 'application/json'}})
  if (res.status == 200) {
    const data = await res.json();
    if (data == null) {
      return {}
    } else {
      Object.keys(data).forEach((k) => {
        data[k] = new Source(k, data[k]);
      });
      return data;
    }
  } else {
    throw new Error((await res.json()).errmsg);
  }
}

export async function getSource(sid) {
  const res = await fetch(`api/sources/${sid}`, {headers: {'Accept': 'application/json'}})
  if (res.status == 200) {
    return new Source(sid, await res.json());
  } else {
    throw new Error((await res.json()).errmsg);
  }
}
