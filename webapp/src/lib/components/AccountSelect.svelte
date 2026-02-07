<script lang="ts">
    import type { EventHandler } from "svelte/elements";
    import { type account } from "../../../types/account";

    let {
        accountArray = $bindable(),
        resolver,
        rejector,
    }: {
        accountArray: account[];
        resolver: ((uuid: string) => void) | null;
        rejector: (() => void) | null;
    } = $props();

    function submitHandler(e: Event) {
        e.preventDefault();
        e.stopPropagation();

        const form = new FormData(e.target as HTMLFormElement);
        console.log("submitting form:", form);
        const uuidField = form.get("account");
        if (!uuidField) {
            console.log("no uuid in account field of form");
            return;
        }
        const uuid = uuidField.valueOf();
        if (typeof uuid != "string") {
            if (rejector) {
                rejector();
            }
            return;
        }
        console.log("uuid:", uuid);
        if (resolver) {
            resolver(uuid);
        }
    }
</script>

<form class="w-fit mx-auto flex flex-col items-center" onsubmit={submitHandler}>
    <h1>Choose an Account:</h1>
    <div class="justify-left flex flex-col">
        {#each accountArray as a}
            <div class="w-fill">
                <label for={a.UUID} class="cursor-pointer hover:text-slate-500">
                    <input
                        id={a.UUID}
                        type="radio"
                        name="account"
                        class="peer hidden"
                        value={a.UUID}
                    />
                    <span>[</span>
                    <span class="invisible peer-checked:visible">x</span>
                    <span>]</span>
                    {a.Name}
                </label>
            </div>
        {/each}
    </div>
    <div class="flex gap-2">
        <button class="cursor-pointer hover:text-slate-500">[ OK ]</button>
        <button
            class="cursor-pointer hover:text-slate-500"
            onclick={() => rejector?.()}>[ CANCEL ]</button
        >
    </div>
</form>
