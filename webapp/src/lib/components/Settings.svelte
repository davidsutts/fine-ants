<script lang="ts">
    import { updateAccounts, type account } from "../../../types/account";

    let { accountArray = $bindable() }: { accountArray: account[] } = $props();

    async function handleAccountsChange(e: Event) {
        accountArray = await updateAccounts(accountArray);
    }
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
                    value={a.Name}
                    onchange={(e: Event) => {
                        const input = e.target as HTMLInputElement;
                        let v = input.value;
                        if (v == "") {
                            accountArray.splice(i, 1);
                            return;
                        }
                        accountArray[i].Name = input.value;
                        handleAccountsChange(e);
                    }}
                /> ]
            </div>
        {/each}
        <div>
            [ <input
                class="text-slate-500"
                type="text"
                placeholder="new account"
                onchange={async (e: Event) => {
                    const input = e.target as HTMLInputElement;
                    accountArray.push({ Name: input.value, UUID: "" });
                    handleAccountsChange(e);
                    input.value = "";
                }}
            /> ]
        </div>
    </div>
</div>
