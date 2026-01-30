import { apiRequest, method } from "../src/lib/api";

export type transaction = {
  EffectiveDate: Date;
  EnteredDate: Date;
  Description: string;
  Amount: number;
  Balance: number;
  Category: string;
};

const routeParseTransactions = "/api/transaction/parse";
const routeCategoriseTransactions = "/api/transaction/categorise";
const routeGetAllTransactions = "/api/transaction/get-all";

export function getAllTransactions(): Promise<transaction[]> {
  return apiRequest<transaction[]>(routeGetAllTransactions, method.GET).then(
    (txs: transaction[]) => {
      return txs.map((tx: transaction) => {
        return {
          EffectiveDate: new Date(tx.EffectiveDate),
          EnteredDate: new Date(tx.EnteredDate),
          Amount: tx.Amount,
          Description: tx.Description,
          Category: tx.Category,
          Balance: tx.Balance,
        };
      });
    },
  );
}

export function parseTransactions(data: FormData) {
  return apiRequest<null>(routeParseTransactions, method.POST, data);
}

export function categoriseTransaction(tx: transaction, category: string) {
  tx.Category = category;
  return apiRequest<transaction>(
    routeCategoriseTransactions,
    method.POST,
    JSON.stringify(tx),
  ).then((tx: transaction) => {
    return {
      EffectiveDate: new Date(tx.EffectiveDate),
      EnteredDate: new Date(tx.EnteredDate),
      Amount: tx.Amount,
      Description: tx.Description,
      Category: tx.Category,
      Balance: tx.Balance,
    };
  });
}
