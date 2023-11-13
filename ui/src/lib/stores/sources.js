import { derived, writable } from 'svelte/store';

import { getSources } from '$lib/source'

function createSourcesStore() {
  const { subscribe, set, update } = writable(null);

  return {
    subscribe,

    set: (v) => {
      update((m) => v);
    },

    refresh: async () => {
      const list = await getSources();
      update((m) => list);
      return list;
    },
  };

}

export const sources = createSourcesStore();

export const sourcesList = derived(
  sources,
  ($sources) => {
    if (!$sources) {
      return [];
    }
    return Object.keys($sources).map((k) => $sources[k]);
  },
);

export const activeSources = derived(
  sourcesList,
  ($sourcesList) => {
    return $sourcesList.filter((s) => s.active);
  },
);
