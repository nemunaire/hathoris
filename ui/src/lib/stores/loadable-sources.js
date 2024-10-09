import { derived, writable } from 'svelte/store';

import { retrieveLoadableSources } from '$lib/custom_source'

function createLoadableSourcesStore() {
  const { subscribe, set, update } = writable(null);

  return {
    subscribe,

    set: (v) => {
      update((m) => v);
    },

    refresh: async () => {
      const map = await retrieveLoadableSources();
      update((m) => map);
      return map;
    },
  };

}

export const loadablesSources = createLoadableSourcesStore();
