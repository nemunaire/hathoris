import { derived, writable } from 'svelte/store';

import { getInputs } from '$lib/input'

function createInputsStore() {
  const { subscribe, set, update } = writable(null);

  return {
    subscribe,

    set: (v) => {
      update((m) => v);
    },

    refresh: async () => {
      const list = await getInputs();
      update((m) => list);
      return list;
    },
  };

}

export const inputs = createInputsStore();

export const inputsList = derived(
  inputs,
  ($inputs) => {
    if (!$inputs) {
      return [];
    }
    return Object.keys($inputs).map((k) => $inputs[k]);
  },
);

export const activeInputs = derived(
  inputsList,
  ($inputsList) => {
    return $inputsList.filter((s) => s.active);
  },
);
