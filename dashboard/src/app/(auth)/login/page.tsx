"use client";
import { fetchLogin } from "@/api/auth";
import Link from "next/link";
import { ChangeEvent, FormEvent, FormEventHandler, useState } from "react";
import { useRouter } from "next/navigation";

function LoginPage() {
  const router = useRouter();

  const [form, setForm] = useState({
    email: "",
    password: ""
  });

  const [error, setError] = useState("");

  function handleFormChange(e: ChangeEvent<HTMLInputElement>) {
    setForm(prev => {
      return {
        ...prev,
        [e.target.name]: e.target.value
      };
    });
  }

  async function submit(e: FormEvent) {
    e.preventDefault();
    const { data, error } = await fetchLogin(form);
    if (error) {
      setError(error);
    } else {
      router.push("/home");
    }
  }
  return (
    <main className="w-[100vw] h-[100vh] flex justify-center items-center">
      <form
        onSubmit={submit}
        className="min-w-[400px] min-h-[600px] p-6 border-4 border-slate-700 rounded-2xl"
      >
        <h1 className="text-4xl my-3 text-center text-slate-700 font-extrabold">
          Login
        </h1>
        <p className="text-center my-2 text-slate-700">
          Click{" "}
          <Link
            className="font-bold underline hover:no-underline"
            href="/register"
          >
            here
          </Link>{" "}
          to register
        </p>
        <p className="text-center my-2 text-red-700">{error}</p>
        <div className="w-full px-4 my-11 flex flex-col">
          <div>
            <label className="label-text" htmlFor="email">
              Email:
            </label>
            <input
              onChange={handleFormChange}
              className="w-full h-8 border-2 border-slate-700 rounded-md focus:ring-4 focus:ring-slate-700"
              type="email"
              name="email"
              id="email"
            />
          </div>
          <div>
            <label className="label-text" htmlFor="password">
              Password:
            </label>
            <input
              onChange={handleFormChange}
              className="w-full h-8 border-2 border-slate-700 rounded-md focus:ring-4 focus:ring-slate-700"
              type="password"
              name="password"
              id="password"
            />
          </div>
          <button className="mt-7 bg-white w-3/4 mx-auto p-1 border-2 border-slate-700 rounded-lg">
            Login
          </button>
        </div>
      </form>
    </main>
  );
}

export default LoginPage;