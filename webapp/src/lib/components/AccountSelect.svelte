<script lang="ts">
    import { type account } from "../../../types/account";

    let { accountArray = $bindable() }: { accountArray: account[] } = $props();
</script>

<form class="w-fit mx-auto flex flex-col items-center">
    <h1>Choose an Account:</h1>
    <div class="justify-left flex flex-col"></div>
    {#each accountArray as a}
        <div class="w-fill">
            <label for={a.UUID}>
                <input
                    id={a.UUID}
                    type="radio"
                    name="account"
                    class="peer hidden"
                />
                <span>[</span>
                <span class="invisible peer-checked:visible">x</span>
                <span>]</span>
                {a.Name}
            </label>
        </div>
    {/each}
</form>
