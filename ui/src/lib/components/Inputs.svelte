<script>
 import { activeInputs, inputs, inputsList } from '$lib/stores/inputs';

 export let showInactives = false;

 let altering_mixer = null;
 async function alterMixer(input, streamid, volume) {
     if (altering_mixer) altering_mixer.abort();
     altering_mixer = setTimeout(() => {
         fetch(`api/inputs/${input.name}/streams/${streamid}/volume`, {headers: {'Accept': 'application/json'}, method: 'POST', body: JSON.stringify({'volume': volume ? volume : input.mixer[streamid].volume})}).then(() => inputs.refresh());
         altering_mixer = null;
     }, 450);
 }
</script>

<ul class="list-group list-group-flush">
    {#if (showInactives && $inputsList.length === 0) || (!showInactives && $activeInputs.length === 0)}
        <li class="list-group-item py-3">
            <span class="text-muted">
                Aucune source active.
            </span>
        </li>
    {/if}
    {#each $inputsList as input, iid}
        {#if showInactives || input.active}
            {#each Object.keys(input.streams) as idstream}
                {@const title = input.streams[idstream]}
                <li class="list-group-item py-3 d-flex flex-column">
                    <div class="d-flex justify-content-between align-items-center">
                        <div>
                            <label for={'input' + iid + 'stream' + idstream} class="form-label d-inline">{title}</label>
                            <span class="text-muted">({input.name})</span>
                        </div>
                        {#if input.controlable}
                            <div>
                                <button
                                    class="btn btn-sm btn-primary"
                                    on:click={() => input.playpause(idstream)}
                                >
                                    <i class="bi bi-pause"></i>
                                </button>
                            </div>
                        {:else if input.mixable && input.mixer[idstream]}
                            <div
                                class="badge bg-primary"
                                title={input.mixer[idstream].volume_percent}
                            >
                                {input.mixer[idstream].volume_db}
                            </div>
                        {/if}
                    </div>
                    {#if input.mixable && input.mixer[idstream]}
                        <div>
                            <input
                                type="range"
                                class="form-range"
                                id={'input' + iid + 'stream' + idstream}
                                min={0}
                                max={65536}
                                bind:value={input.mixer[idstream].volume}
                                on:change={() => alterMixer(input, idstream)}
                            >
                        </div>
                    {/if}
                </li>
            {/each}
        {/if}
    {/each}
</ul>
