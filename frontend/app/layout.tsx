import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "EGOT Tracker",
  description: "Track celebrity progress toward EGOT status - Emmy, Grammy, Oscar, Tony",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-hollywood-black font-body antialiased">
        <div className="spotlight min-h-screen">
          {children}
        </div>
      </body>
    </html>
  );
}
