<script lang="ts">
    let { accountArray = $bindable() }: { accountArray: string[] } = $props();
</script>

<div class="border border-dashed my-5"></div>
<div class="mx-auto w-fit grid grid-cols-2 gap-3">
    <h3 class="text-right">[ ACCOUNTS ] :</h3>
    <div class=" flex flex-col gap-2">
        {#each accountArray as a, i}
            <div>
                [ <input
                    class="text-slate-500"
                    type="text"
                    value={a}
                    onchange={(e: Event) => {
                        const input = e.target as HTMLInputElement;
                        let v = input.value;
                        if (v == "") {
                            accountArray.splice(i, 1);
                            return;
                        }
                        accountArray[i] = input.value;
                    }}
                /> ]
            </div>
        {/each}
        <div>
            [ <input
                class="text-slate-500"
                type="text"
                placeholder="new account"
                onchange={(e: Event) => {
                    const input = e.target as HTMLInputElement;
                    accountArray.push(input.value);
                    input.value = "";
                }}
            /> ]
        </div>
    </div>
</div>
