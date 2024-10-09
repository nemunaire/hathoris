import { error } from '@sveltejs/kit';
import { getSettings } from '$lib/settings';
import { loadablesSources } from '$lib/stores/loadable-sources.js';

export const load = async({ parent, fetch }) => {
  const data = await parent();

  await loadablesSources.refresh();

  let settings;
  try {
    settings = await getSettings();
  } catch (err) {
    throw error(err.status, err.statusText);
  }

  const version = fetch('/api/version').then((res) => res.json());

  return {
    ...data,
    settings,
    version,
  };
}
