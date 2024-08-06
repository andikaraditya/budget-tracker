"use server";

import { apiResponse, getToken } from "./api";

const BASE_URL = process.env.BACKEND_HOST;

export async function fetchSummary() {
  const token = getToken();
  const res = await fetch(`${BASE_URL}/summary`, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${token}`
    },
    cache: "no-store"
  });

  const result = await res.json();

  let response: apiResponse;
  if (res.ok) {
    response = {
      status: res.status,
      data: result,
      error: null
    };
  } else {
    response = {
      status: res.status,
      data: null,
      error: result.error
    };
  }
  return response;
}
