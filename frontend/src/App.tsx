import { useState } from "react";
import Sidebar from "./components/Sidebar";
import Dashboard from "./components/Dashboard";
import ConfigEditor from "./components/ConfigEditor";

export default function App() {
  const [activeTab, setActiveTab] = useState("dashboard");
  const [activeConfig, setActiveConfig] = useState("");

  return (
    <div className="flex h-screen w-screen bg-slate-900 text-white overflow-hidden font-sans">
      <Sidebar
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        activeConfig={activeConfig}
        setActiveConfig={setActiveConfig}
      />

      <main className="flex-1 h-full min-w-0 bg-slate-950/50">
        {activeTab === "dashboard" && <Dashboard activeConfig={activeConfig} />}

        {activeTab === "editor" && <ConfigEditor activeConfig={activeConfig} />}
      </main>
    </div>
  );
}
