import NavBar from "@/components/NavBar";

function DashboardLayout({
  children
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <>
      <NavBar />
      {children}
    </>
  );
}

export default DashboardLayout;
