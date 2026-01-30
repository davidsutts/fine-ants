<script lang="ts">
    import upload from "./assets/upload.svg";
    import {
        getAllTransactions,
        parseTransactions,
        type transaction,
    } from "../types/transaction";
    import { onMount } from "svelte";

    let transactions: transaction[] = $state([]);

    const modeList = "List";
    const modeSort = "Sort";
    let pageMode: number = $state(1);
    const modes = [modeList, modeSort];

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
    ];

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
    }

    onMount(async () => {
        transactions = await getAllTransactions();
    });
</script>

{#snippet element()}
    {#if modes[pageMode] == modeList}
        {@render table()}
    {:else if modes[pageMode] == modeSort}
        {@render sort()}
    {/if}
{/snippet}

{#snippet table()}
    <div class="grid grid-cols-6 shrink m-auto">
        {#each transactions as transaction}
            <p class="col-span-1">
                {transaction.EnteredDate.toLocaleDateString("en-AU")}
            </p>
            <p class="col-span-1 font-bold text-lg text-right pr-5">
                {money.format(transaction.Amount)}
            </p>
            <p class="col-span-3">{transaction.Description}</p>
            <div class="col-span-1 flex gap-2">
                <p class="text-gray-500">BAL:</p>
                {transaction.Balance}
            </div>
        {/each}
    </div>
{/snippet}

{#snippet categoryButton(c: category, index: number)}
    <button
        class="border px-2 col-span-1 cursor-pointer hover:border-gray-400 hover:text-gray-400"
    >
        <p class="text-gray-500">{index}</p>
        <p>{c.emoji}</p>
        <p>{c.name}</p>
    </button>
{/snippet}

{#snippet sort()}
    {#if transactions.length > 0}
        <div class="m-auto flex-col gap-5 flex">
            <div class="border p-5 flex flex-col items-center gap-3">
                <p class="w-full text-left">
                    {transactions[0].EffectiveDate.toLocaleDateString()}
                </p>
                <p class="text-5xl">
                    {money.format(transactions[0].Amount)}
                </p>
                <p>{transactions[0].Description}</p>
            </div>
            <div class="grid grid-cols-5 gap-5">
                {#each categories as c, i}
                    {@render categoryButton(c, i + 1)}
                {/each}
            </div>
        </div>
    {:else}
        <div class="border p-5">Upload transactions to get started.</div>
    {/if}
{/snippet}

<main class="p-5 font-mono">
    <header class="mb-5 flex gap-5">
        <div>
            <h1 class="text-2xl">🐜 Fine-Ants</h1>
            <p>Finance Done Properly</p>
        </div>
        <label
            for="avatar"
            class="hover:cursor-pointer w-24 h-full shrink-0 border px-2 text-center"
        >
            <strong>Upload</strong>
            <p>↥</p>
        </label>
        <button
            class="hover:cursor-pointer border px-2 h-full w-24"
            onmousedown={() => (pageMode = (pageMode + 1) % 2)}
        >
            <strong>Mode</strong>
            <p>{modes[pageMode]}</p>
        </button>
    </header>
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
</main>

<style>
</style>
