<script>
 function refreshMixers() {
     const mxrs = fetch('api/mixer', {headers: {'Accept': 'application/json'}}).then((res) => res.json());
     mxrs.then((m) => {
         mixers = m;
         altering_mixer = null;
     })
 }

 export let showReadOnly = false;
 export let advanced = false;

 let mixers = null;
 refreshMixers();
 setInterval(refreshMixers, 5000);

 let altering_mixer = null;
 async function alterMixer(mixer, values) {
     if (altering_mixer) altering_mixer.abort();
     altering_mixer = setTimeout(() => {
         fetch(`api/mixer/${mixer.NumID}/values`, {headers: {'Accept': 'application/json'}, method: 'POST', body: JSON.stringify(values ? values : (advanced ? mixer.values : [mixer.values[0]]))}).then(refreshMixers);
         altering_mixer = null;
     }, 450);
 }
</script>

{#if !mixers}
    <div class="card-body d-flex flex-fill justify-content-center">
        <div class="spinner-border" role="status">
            <span class="visually-hidden">Loading...</span>
        </div>
    </div>
{:else}
    <ul class="list-group list-group-flush">
    {#each mixers as mixer (mixer.NumID)}
        {#if showReadOnly || mixer.RW}
            <li class="list-group-item py-3">
                {#if mixer.items}
                    <label for={mixer.Name + '0'} class="form-label">{mixer.Name}</label>
                    {#if mixer.values}
                        {#each mixer.values as cur, idx}
                            <select
                                class="form-select"
                                disabled={!mixer.RW}
                                id={mixer.Name + idx}
                                bind:value={cur}
                                on:change={() => alterMixer(mixer)}
                            >
                                {#each mixer.items as opt, idx}
                                    <option value={idx}>{opt}</option>
                                {/each}
                            </select>
                        {/each}
                    {/if}
                {:else if mixer.Type === "INTEGER"}
                    <label for={mixer.Name + '0'} class="form-label">{mixer.Name}</label>
                    {#if mixer.values}
                        <div class="badge bg-primary float-end">
                            {mixer.DBScale ? ((mixer.DBScale.Min + mixer.DBScale.Step * mixer.values[0]) + ' dB') : mixer.values[0]}
                        </div>
                        {#each mixer.values as cur, idx}
                            {#if advanced || idx === 0}
                                <input
                                    type="range"
                                    class="form-range"
                                    disabled={!mixer.RW}
                                    id={mixer.Name + idx}
                                    min={mixer.Min}
                                    max={mixer.Max}
                                    step={mixer.Step}
                                    title={mixer.DBScale ? ((mixer.DBScale.Min + mixer.DBScale.Step * cur) + ' dB') : cur}
                                    bind:value={cur}
                                    on:change={() => alterMixer(mixer)}
                                >
                            {/if}
                        {/each}
                    {/if}
                {:else if mixer.Type === "BOOLEAN"}
                    {#if mixer.RW}
                        <div class="btn-group" role="group" aria-label="Basic example">
                            <button
                                class="btn"
                                class:btn-secondary={!mixer.values.reduce((a,b) => a || b)}
                                class:btn-primary={mixer.values.reduce((a,b) => a || b)}
                                on:click={() => alterMixer(mixer, [!mixer.values.reduce((a,b) => a || b)])}
                            >
                                {mixer.Name}
                            </button>
                            {#if advanced && mixer.values.length > 1}
                                {#each mixer.values as cur, ichan}
                                    <button
                                        class="btn btn-sm"
                                        class:btn-secondary={!cur}
                                        class:btn-primary={cur}
                                        on:click={() => alterMixer(mixer, mixer.values.map((v, i) => i == ichan ? !v : v))}
                                    >
                                        {ichan+1}
                                    </button>
                                {/each}
                            {/if}
                        </div>
                    {:else}
                        <label class="form-label">{mixer.Name}</label>
                    {/if}
                {/if}
            </li>
        {/if}
    {/each}
    </ul>
{/if}
