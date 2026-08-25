import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "[[.ProjectName]] | Powered by Koko",
  description: "Modern fullstack application built with Turborepo and Koko-cli",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="antialiased min-h-screen bg-[#09090b] text-[#fafafa]">
        {children}
      </body>
    </html>
  );
}
