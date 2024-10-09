export class CustomSource {
  constructor(res) {
    if (res) {
      this.update(res);
    }
  }

  update({ src, kv }) {
    this.src = src;
    this.kv = kv;
  }
}

export async function retrieveLoadableSources() {
  const res = await fetch(`api/settings/loadable_sources`, {headers: {'Accept': 'application/json'}})
  if (res.status == 200) {
    const data = await res.json();
    if (data == null) {
      return {}
    } else {
      return data;
    }
  } else {
    throw new Error((await res.json()).errmsg);
  }
}
