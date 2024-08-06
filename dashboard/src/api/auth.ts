"use server";
import { cookies } from "next/headers";
import { apiResponse } from "./api";
import { redirect } from "next/navigation";

type LoginInput = {
  email: string;
  password: string;
};

type RegisterInput = {
  email: string;
  password: string;
};

const BASE_URL = process.env.BACKEND_HOST;

export async function fetchLogin(input: LoginInput) {
  const cookieStore = cookies();
  const res = await fetch(`${BASE_URL}/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(input),
    cache: "no-store"
  });
  const result = await res.json();
  let response: apiResponse;
  if (res.ok) {
    response = {
      status: res.status,
      data: "success",
      error: null
    };
    const oneDay = 1 * 24 * 60 * 60;

    cookieStore.set("token", result.token, { maxAge: oneDay });
  } else {
    response = {
      status: res.status,
      data: null,
      error: result.error
    };
  }
  return response;
}

export async function fetchRegister(input: RegisterInput) {
  const res = await fetch(`${BASE_URL}/register`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(input),
    cache: "no-store"
  });
  const result = await res.json();
  let response: apiResponse;
  if (res.ok) {
    response = {
      status: res.status,
      data: "success",
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

export async function handleLogout() {
  const cookieStore = cookies();
  cookieStore.delete("token");
  redirect("/login");
}
