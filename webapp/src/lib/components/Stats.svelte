<script lang="ts">
    import { type account } from "../../../types/account";
    import { type category } from "../../../types/category";
    import type { transaction } from "../../../types/transaction";

    let {
        accountArray = $bindable(),
        transactions = $bindable(),
        categories,
    }: {
        accountArray: account[];
        transactions: transaction[];
        categories: category[];
    } = $props();

    type catValue = {
        dollars: number;
        normalised: number;
    };

    let processedVals = $derived.by(() => {
        let buckets = new Map<string, catValue>();
        let maxAbsValue = 0; // Find the "peak" value for scaling
        transactions.forEach((tx) => {
            if (tx.Category == "") {
                return;
            }
            let cur = buckets.get(tx.Category);
            if (cur == undefined) {
                cur = { dollars: 0, normalised: 0 };
                buckets.set(tx.Category, cur);
            }
            cur.dollars += tx.Amount;
        });

        buckets.forEach((val) => {
            const abs = Math.abs(val.dollars);
            if (abs > maxAbsValue) {
                maxAbsValue = abs;
                console.log("found a bigger val:", maxAbsValue);
            }
        });

        buckets.forEach((val) => {
            val.normalised = maxAbsValue > 0 ? val.dollars / maxAbsValue : 0;
        });
        return buckets;
    });

    const today = new Date();
    const thisYear = today.getFullYear();
    const years = new Array(3).fill(0).map((_, i) => {
        return thisYear - i;
    });
    console.info(years);

    const money = Intl.NumberFormat("en-AU", {
        currency: "AUD",
        style: "currency",
    });

    function printCategoryDollars(cat: string): string {
        let v = processedVals.get(cat);
        if (v == undefined) {
            v = { dollars: 0, normalised: 0 };
        }

        return money.format(v.dollars);
    }

    function getNormalisedVal(cat: string): number {
        let v = processedVals.get(cat);
        if (v == undefined) return 0;
        return v.normalised * 75; // Leave room for
    }
</script>

{#snippet filters()}
    <h2 class="font-bold text-lg">Filters</h2>
    <h3 class="font-bold">Account:</h3>
    <div class="flex gap-5">
        {#each accountArray as acc}
            <div class="flex">
                <label
                    for="acc-{acc.Name}"
                    class="cursor-pointer hover:text-slate-500"
                >
                    <input
                        type="checkbox"
                        id="acc-{acc.Name}"
                        class="peer hidden"
                    />
                    <span>[</span><span class="invisible peer-checked:visible"
                        >x</span
                    ><span>]</span>
                    {acc.Name}
                </label>
            </div>
        {/each}
    </div>
    <h3 class="font-bold">Year:</h3>
    <div class="flex gap-5">
        {#each years as year}
            <label class="cursor-pointer hover:text-slate-500">
                <input type="checkbox" class="peer hidden" />
                <span>[</span><span class="invisible peer-checked:visible"
                    >x</span
                ><span>]</span>
                {year}
            </label>
        {/each}
    </div>
{/snippet}

{#snippet categoryBar(cat: category)}
    {@const norm = getNormalisedVal(cat.name)}
    {@const absoluteWidth = Math.abs(norm)}

    <div class="flex items-center justify-between w-full gap-4 py-4">
        <p class="whitespace-nowrap shrink-0 w-48">
            {cat.emoji} <span class="ml-2">{cat.name}</span>
        </p>

        <div class="flex-grow flex flex-row-reverse items-center relative h-5">
            <div
                class="h-full transition-all duration-500 shrink-0"
                class:bg-red-500={norm < 0}
                class:bg-slate-700={norm >= 0}
                style="width: {absoluteWidth}%"
            ></div>

            <div
                class="pr-4 tabular-nums tracking-wide"
                class:text-red-500={norm < 0}
            >
                {printCategoryDollars(cat.name)}
            </div>
        </div>
    </div>
{/snippet}

<div class="select-none w-full">
    <!-- {@render filters()} -->
    <div class="w-fill">
        {#each categories as cat}
            {@render categoryBar(cat)}
        {/each}
    </div>
</div>
