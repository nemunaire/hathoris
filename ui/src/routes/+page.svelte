<script>
 import Applications from '$lib/components/Applications.svelte';
 import Inputs from '$lib/components/Inputs.svelte';
 import Mixer from '$lib/components/Mixer.svelte';
 import { activeSources } from '$lib/stores/sources';
 import { activeInputs } from '$lib/stores/inputs';

 let mixerAdvanced = false;
 let showInactiveInputs = false;
</script>

<div class="container">
    <div class="row">
        <div class="col-md">
            <div class="card my-3">
                <h4 class="card-header">
                    <div class="d-flex justify-content-between">
                        <div>
                            <i class="bi bi-sliders"></i>
                            Mixer
                        </div>
                        <button
                            class="btn btn-sm"
                            class:btn-info={mixerAdvanced}
                            class:btn-secondary={!mixerAdvanced}
                            on:click={() => { mixerAdvanced = !mixerAdvanced; }}
                        >
                            <i class="bi bi-alt"></i>
                        </button>
                    </div>
                </h4>
                <Mixer advanced={mixerAdvanced} />
            </div>
        </div>

        <div class="col-md">
            {#if $activeSources.length > 0}
                <div class="card my-3">
                    <h4 class="card-header">
                        <div class="d-flex justify-content-between">
                            <div>
                                <i class="bi bi-window-stack"></i>
                                Applications
                            </div>
                        </div>
                    </h4>
                    <Applications />
                </div>
            {/if}
            <div class="card my-3">
                <h4 class="card-header">
                    <div class="d-flex justify-content-between">
                        <div>
                            <i class="bi bi-speaker"></i>
                            Sources
                        </div>
                        <button
                            class="btn btn-sm"
                            class:btn-info={showInactiveInputs}
                            class:btn-secondary={!showInactiveInputs}
                            on:click={() => { showInactiveInputs = !showInactiveInputs; }}
                        >
                            <i class="bi bi-eye-fill"></i>
                        </button>
                    </div>
                </h4>
                <Inputs showInactive={showInactiveInputs} />
            </div>
        </div>
    </div>
</div>
