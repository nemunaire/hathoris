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
 async function muteMixer(input, streamid, mute) {
     fetch(`api/inputs/${input.name}/streams/${streamid}/volume`, {headers: {'Accept': 'application/json'}, method: 'POST', body: JSON.stringify({'mute': mute !== undefined ? mute : input.mixer[streamid].mute})}).then(() => inputs.refresh());
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
                        <div class="d-flex align-items-center">
                            {#if input.mixable && input.mixer[idstream]}
                                <button
                                    class="btn btn-sm ms-1"
                                    class:btn-primary={input.mixer[idstream].mute}
                                    class:btn-secondary={!input.mixer[idstream].mute}
                                    on:click={() => {muteMixer(input, idstream, !input.mixer[idstream].mute);}}
                                >
                                    <i class="bi bi-volume-mute-fill"></i>
                                </button>
                            {/if}
                            {#if input.controlable}
                                <button
                                    class="btn btn-sm btn-primary ms-1"
                                    on:click={() => input.playpause(idstream)}
                                >
                                    <i class="bi bi-pause"></i>
                                </button>
                            {:else if input.mixable && input.mixer[idstream]}
                                <div
                                    class="badge bg-primary ms-1"
                                    title={input.mixer[idstream].volume_percent}
                                >
                                    {input.mixer[idstream].volume_db}
                                </div>
                            {/if}
                        </div>
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
