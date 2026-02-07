import { apiRequest, method } from "../src/lib/api";

export type account = {
  UUID: string;
  Name: string;
};

const routeUpdateAccount = "/api/accounts/update";
const routeGetAllAccounts = "/api/accounts/get-all";

export function updateAccounts(accs: account[]): Promise<account[]> {
  return apiRequest<account[]>(
    routeUpdateAccount,
    method.POST,
    JSON.stringify(accs),
  );
}

export function getAllAccounts(): Promise<account[]> {
  return apiRequest<account[]>(routeGetAllAccounts, method.GET);
}
