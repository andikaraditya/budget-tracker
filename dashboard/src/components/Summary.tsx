"use client";

import { fetchSummary } from "@/api/summary";
import { useEffect, useState } from "react";

function Summary() {
  const [summary, setSummary] = useState({
    expense: 0,
    income: 0,
    total: 0
  });
  useEffect(() => {
    fetchSummary()
      .then(result => {
        const { status, data, error } = result;
        if (error) {
          throw error;
        }
        setSummary({
          ...data
        });
      })
      .catch(err => {
        console.log(err);
      });
  }, []);

  return (
    <>
      <div className="w-[400px] text-slate-900 border-b-4 border-slate-700 mt-6 mb-2">
        <h1 className="font-bold text-center text-4xl mb-4">Summary</h1>
      </div>
      <div className="w-[600px] h-[100px] flex flex-col justify-center gap-4">
        <div className="w-full flex justify-around">
          <p className="flex-1 text-center font-semibold text-2xl ">Expense</p>
          <p className="flex-1 text-center font-semibold text-2xl">Income</p>
          <p className="flex-1 text-center font-semibold text-2xl">Total</p>
        </div>
        <div className="w-full flex justify-around">
          <p className="flex-1 text-center font-semibold text-xl text-rose-600">
            {summary.expense}
          </p>
          <p className="flex-1 text-center font-semibold text-xl text-green-600">
            {summary.income}
          </p>
          <p
            className={`flex-1 text-center font-semibold text-xl ${
              summary.total > 0 ? "text-green-600" : "text-rose-600"
            }`}
          >
            {summary.total}
          </p>
        </div>
      </div>
    </>
  );
}

export default Summary;
