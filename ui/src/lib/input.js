export class Input {
  constructor(id, res) {
    this.id = id;
    if (res) {
      this.update(res);
    }
  }

  update({ name, active, controlable, hasplaylist, streams, mixable, mixer }) {
    this.name = name;
    this.active = active;
    this.controlable = controlable;
    this.hasplaylist = hasplaylist;
    this.streams = streams;
    this.mixable = mixable;
    this.mixer = mixer;
  }

  async getStreams() {
    const data = await fetch(`api/inputs/${this.id}/streams`, {headers: {'Accept': 'application/json'}});
    if (data.status == 200) {
      return await data.json();
    } else {
      throw new Error((await res.json()).errmsg);
    }
  }

  async playpause(idstream) {
    const data = await fetch(`api/inputs/${this.id}/streams/${idstream}/pause`, {headers: {'Accept': 'application/json'}, method: 'POST'});
    if (data.status != 200) {
      throw new Error((await res.json()).errmsg);
    }
  }

  async nexttrack(idstream) {
    const data = await fetch(`api/inputs/${this.id}/streams/${idstream}/next_track`, {headers: {'Accept': 'application/json'}, method: 'POST'});
    if (data.status != 200) {
      throw new Error((await res.json()).errmsg);
    }
  }

  async nextrandomtrack(idstream) {
    const data = await fetch(`api/inputs/${this.id}/streams/${idstream}/next_random_track`, {headers: {'Accept': 'application/json'}, method: 'POST'});
    if (data.status != 200) {
      throw new Error((await res.json()).errmsg);
    }
  }

  async prevtrack(idstream) {
    const data = await fetch(`api/inputs/${this.id}/streams/${idstream}/prev_track`, {headers: {'Accept': 'application/json'}, method: 'POST'});
    if (data.status != 200) {
      throw new Error((await res.json()).errmsg);
    }
  }
}

export async function getInputs() {
  const res = await fetch(`api/inputs`, {headers: {'Accept': 'application/json'}})
  if (res.status == 200) {
    const data = await res.json();
    if (data == null) {
      return {}
    } else {
      Object.keys(data).forEach((k) => {
        data[k] = new Input(k, data[k]);
      });
      return data;
    }
  } else {
    throw new Error((await res.json()).errmsg);
  }
}

export async function getInput(sid) {
  const res = await fetch(`api/inputs/${sid}`, {headers: {'Accept': 'application/json'}})
  if (res.status == 200) {
    return new Input(sid, await res.json());
  } else {
    throw new Error((await res.json()).errmsg);
  }
}
