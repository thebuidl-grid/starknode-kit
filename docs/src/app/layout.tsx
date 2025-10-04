import type { Metadata } from "next";
import { Inter, Source_Code_Pro } from "next/font/google";
import "./globals.css";
import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
});

const source_code_pro = Source_Code_Pro({
  subsets: ["latin"],
  variable: "--font-source-code"
})

export const metadata: Metadata = {
  title: "starknode-kit Documentation",
  description: "Complete documentation for starknode-kit - A CLI tool for setting up and managing Ethereum and Starknet nodes",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${inter.variable} ${source_code_pro.variable}`}>
      <body className="antialiased bg-white text-gray-900">
        <Sidebar />
        <Header />
        <main className="ml-64 pt-16 min-h-screen bg-white">
          <div className="max-w-5xl mx-auto px-12 py-16">
            {children}
          </div>
        </main>
      </body>
    </html>
  );
}
