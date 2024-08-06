import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import LogoutButton from "./LogoutButton";

export default function NavBar() {
  const cookieStore = cookies();
  return (
    <nav className="flex justify-between px-7 py-5 border-b-4 border-slate-700">
      <div className="flex-1"></div>
      <div className="flex-1">
        <p className="text-2xl text-center font-bold text-slate-700">
          Budget Tracker
        </p>
      </div>
      <div className="flex-1 flex flex-row-reverse">
        <LogoutButton />
      </div>
    </nav>
  );
}
