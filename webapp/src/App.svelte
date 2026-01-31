<script lang="ts">
    import {
        categoriseTransaction,
        getAllTransactions,
        parseTransactions,
        type transaction,
    } from "../types/transaction";

    import { onMount } from "svelte";
    import Settings from "./lib/components/Settings.svelte";

    let transactions: transaction[] = $state([]);
    let sortIndex: number = $state(0);

    const modeList = "LIST";
    const modeSort = "SORT";
    let pageMode: number = $state(0);
    const modes = [modeList, modeSort];

    $effect(() => {
        $inspect(accounts);
    });

    let accounts: string[] = $state(["Transactions"]);
    let showSettings: boolean = $state(true);

    type category = {
        name: string;
        emoji: string;
    };

    const money = Intl.NumberFormat("en-AU", {
        currency: "AUD",
        style: "currency",
    });

    const categories: category[] = [
        { name: "Work", emoji: "💼" },
        { name: "Transport", emoji: "🚗" },
        { name: "House", emoji: "🏠" },
        { name: "Groceries", emoji: "🥬" },
        { name: "Utilities", emoji: "⚡" },

        { name: "Insurance", emoji: "🗄️" },
        { name: "Health", emoji: "💊" },
        { name: "Sport", emoji: "🏃" },
        { name: "Donations", emoji: "🙌" },
        { name: "Personal", emoji: "💰" },

        { name: "Dining", emoji: "☕" },
        { name: "Gift", emoji: "🎁" },
        { name: "Transfer", emoji: "↔️" },
        { name: "Skip", emoji: "⏩" },
        { name: "Back", emoji: "⏪" },
    ];

    function getNewSortIndex() {
        for (let i = 0; i < transactions.length; i++) {
            console.log(transactions[i].Category);
            if (transactions[i].Category == "") {
                sortIndex = i;
                return;
            }
            if (i == transactions.length - 1) {
                // We haven't found an empty transaction.
                sortIndex = -1;
                return;
            }
        }
    }

    async function handleFileUpload(e: Event) {
        const input = e.target as HTMLInputElement;

        let files = input.files;
        if (files == null || files.length <= 0) {
            return;
        }

        let data = new FormData();
        data.append("file", files[0], files[0].name);

        await parseTransactions(data);
        transactions = await getAllTransactions();
        getNewSortIndex();
    }

    onMount(async () => {
        transactions = await getAllTransactions();
        getNewSortIndex();
    });

    async function onCategorise(c: category) {
        if (c.name == "Back") {
            sortIndex--;
            return;
        }
        if (c.name != "Skip") {
            transactions[sortIndex] = await categoriseTransaction(
                transactions[sortIndex],
                c.name,
            );
        }
        sortIndex++;
        while (transactions[sortIndex].Category != "") {
            sortIndex++;
            console.log(
                "transaction[%d] = %s",
                sortIndex,
                transactions[sortIndex].Category,
            );
        }
    }
</script>

{#snippet element()}
    {#if modes[pageMode] == modeList}
        {@render table()}
    {:else if modes[pageMode] == modeSort}
        {@render sort()}
    {/if}
{/snippet}

{#snippet table()}
    <div class="grid grid-cols-6 shrink mx-auto gap-5">
        {#each transactions as transaction}
            <p class="col-span-1">
                {transaction.EnteredDate.toLocaleDateString("en-AU")}
            </p>
            <p class="col-span-1 font-bold text-lg text-right pr-5"></p>
            <div class="col-span-3 flex flex-col">
                <p>{transaction.Description}</p>
                <p class="text-slate-500">{transaction.Category}</p>
            </div>
            <div class="col-span-1 flex flex-col">
                <p class="font-bold text-right">
                    {money.format(transaction.Amount)}
                </p>
                <div class="flex flex-row justify-between text-slate-500">
                    <p>BAL:</p>
                    {money.format(transaction.Balance)}
                </div>
            </div>
        {/each}
    </div>
{/snippet}

{#snippet categoryButton(c: category, index: number)}
    <button
        class="px-2 col-span-1 cursor-pointer hover:border-gray-400 hover:text-gray-400"
        onclick={() => onCategorise(c)}
    >
        <p class="text-slate-500">{index}</p>
        <p>{c.emoji}</p>
        <p>[ {c.name} ]</p>
    </button>
{/snippet}

{#snippet sort()}
    {#if transactions.length > 0 && sortIndex >= 0}
        <div class="m-auto flex-col gap-5 flex">
            <div class="p-5 flex flex-col items-center gap-3">
                <div class="flex justify-between w-full">
                    <p>
                        {transactions[
                            sortIndex
                        ].EffectiveDate.toLocaleDateString()}
                    </p>
                    <p>{transactions[sortIndex].Category}</p>
                </div>
                <p class="text-5xl">
                    {money.format(transactions[sortIndex].Amount)}
                </p>
                <p>[ {transactions[sortIndex].Description} ]</p>
            </div>
            <div class="w-full border border-dashed my-5"></div>
            <div class="grid grid-cols-5 gap-5">
                {#each categories as c, i}
                    {@render categoryButton(c, i + 1)}
                {/each}
            </div>
        </div>
    {:else if transactions.length == 0}
        <div class="mx-auto p-5">Upload transactions to get started.</div>
    {:else if sortIndex <= 0}
        <div class="p-5 mx-auto">[ All transactions have been sorted ]</div>
    {/if}
{/snippet}

<main class="p-5 font-[Cutive_Mono] bg-slate-600 min-h-screen w-screen">
    <div class="bg-slate-200 p-10 w-fit mx-auto">
        <header class="mb-5 flex mx-auto w-fit items-center flex-col">
            <h1 class="text-4xl w-fit font-bold font-sans">🐜 Fine-Ants</h1>
            <p class="text-xl w-fit text-slate-500">Finance Done Properly</p>
        </header>
        <div class="flex gap-5 mx-auto w-fit">
            <label
                for="avatar"
                class="hover:cursor-pointer hover:text-slate-500 text-lg"
            >
                [ UPLOAD ↥ ]
            </label>
            <a
                class="hover:cursor-pointer hover:text-slate-500 text-lg"
                href="/api/transaction/download"
            >
                [ DOWNLOAD ↧ ]
            </a>
            <button
                class="hover:cursor-pointer hover:text-slate-500 text-lg"
                onmousedown={() => (pageMode = (pageMode + 1) % 2)}
            >
                [ {modes[pageMode]} ]
            </button>
            <button
                class="hover:cursor-pointer hover:text-slate-500 text-lg"
                onclick={() => {
                    showSettings = !showSettings;
                }}>[ SETTINGS ]</button
            >
        </div>
        {#if showSettings}
            <Settings bind:accountArray={accounts}></Settings>
        {/if}
        <div class="w-full border border-dashed my-5"></div>
        <input
            accept="text/csv"
            id="avatar"
            type="file"
            class="hidden"
            onchange={handleFileUpload}
        />
        <div class="flex gap-5">
            {@render element()}
        </div>
    </div>
</main>

<style>
    @import url("https://fonts.googleapis.com/css2?family=Cutive+Mono&display=swap");
</style>
