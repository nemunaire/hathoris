<script>
 import { page } from '$app/stores';

 export { className as class };
 let className = '';

 function showModal() {

 }
</script>

<div class={className}>
    <a
        href={$page.route.id == '/settings' ? "/" : "/settings"}
        class="btn btn-link"
        class:text-muted={$page.route.id == '/settings'}
    >
        <i class="bi bi-gear"></i>
    </a>
</div>
