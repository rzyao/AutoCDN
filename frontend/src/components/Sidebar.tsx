import { useState, useEffect } from "react";
import {
  Settings,
  Activity,
  FileText,
  Plus,
  Database,
  Check,
  X,
  Trash2,
} from "lucide-react";
import {
  GetConfigList,
  SaveConfig,
  CreateNewConfig,
  DeleteConfig,
  LoadConfig,
} from "../../wailsjs/go/main/App";
import clsx from "clsx";
import { config } from "../../wailsjs/go/models";

interface SidebarProps {
  activeTab: string;
  setActiveTab: (tab: string) => void;
  activeConfig: string;
  setActiveConfig: (config: string) => void;
}

export default function Sidebar({
  activeTab,
  setActiveTab,
  activeConfig,
  setActiveConfig,
}: SidebarProps) {
  const [configs, setConfigs] = useState<string[]>([]);
  const [isCreating, setIsCreating] = useState(false);
  const [newConfigName, setNewConfigName] = useState("");

  useEffect(() => {
    refreshConfigs();
  }, []);

  const refreshConfigs = async () => {
    try {
      const list = await GetConfigList();
      setConfigs(list || []);
      if (list && list.length > 0 && !activeConfig) {
        setActiveConfig(list[0]);
      }
    } catch (e) {
      console.error(e);
    }
  };

  const handleCreateConfig = async () => {
    if (!newConfigName.trim()) return;

    let name = newConfigName.trim();
    if (!name.endsWith(".yaml") && !name.endsWith(".yml")) {
      name += ".yaml";
    }

    try {
      // Create config with defaults
      await CreateNewConfig(name);
      await refreshConfigs();
      setActiveConfig(name);
      setIsCreating(false);
      setNewConfigName("");
    } catch (e) {
      console.error(e);
      alert("创建失败: " + e);
    }
  };

  return (
    <div className="w-64 h-full bg-slate-900/50 backdrop-blur-xl border-r border-white/10 flex flex-col">
      <div className="p-6 flex items-center gap-3 border-b border-white/5">
        <div className="w-8 h-8 bg-blue-500 rounded-lg flex items-center justify-center">
          <Activity className="text-white w-5 h-5" />
        </div>
        <h1 className="text-xl font-bold bg-linear-to-r from-blue-400 to-cyan-400 bg-clip-text text-transparent">
          AutoCDN
        </h1>
      </div>

      <div className="flex-1 overflow-y-auto py-4">
        <div className="px-4 mb-2 text-xs font-medium text-slate-400 uppercase tracking-wider">
          菜单
        </div>
        <nav className="space-y-1 px-2">
          <button
            onClick={() => setActiveTab("dashboard")}
            className={clsx(
              "w-full flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors",
              activeTab === "dashboard"
                ? "bg-blue-500/10 text-blue-400"
                : "text-slate-400 hover:bg-white/5 hover:text-slate-200",
            )}
          >
            <Activity className="w-4 h-4" />
            仪表盘
          </button>
          <button
            onClick={() => setActiveTab("editor")}
            className={clsx(
              "w-full flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors",
              activeTab === "editor"
                ? "bg-blue-500/10 text-blue-400"
                : "text-slate-400 hover:bg-white/5 hover:text-slate-200",
            )}
          >
            <Settings className="w-4 h-4" />
            配置编辑
          </button>
        </nav>

        <div className="mt-8 px-4 mb-2 text-xs font-medium text-slate-400 uppercase tracking-wider flex justify-between items-center">
          <span>配置文件</span>
          <button
            onClick={() => setIsCreating(true)}
            className="hover:text-white transition-colors p-1"
            title="新建配置"
          >
            <Plus className="w-3.5 h-3.5" />
          </button>
        </div>

        {isCreating && (
          <div className="px-2 mb-2">
            <div className="flex items-center gap-1 bg-slate-800 rounded-md p-1 border border-blue-500/50">
              <input
                autoFocus
                value={newConfigName}
                onChange={(e) => setNewConfigName(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && handleCreateConfig()}
                placeholder="config_new..."
                className="w-full bg-transparent text-xs text-white px-2 focus:outline-none"
              />
              <button
                onClick={handleCreateConfig}
                className="text-emerald-400 hover:bg-white/10 p-1 rounded"
              >
                <Check className="w-3 h-3" />
              </button>
              <button
                onClick={() => setIsCreating(false)}
                className="text-red-400 hover:bg-white/10 p-1 rounded"
              >
                <X className="w-3 h-3" />
              </button>
            </div>
          </div>
        )}

        <div className="space-y-1 px-2">
          {configs.map((cfg) => (
            <div
              key={cfg}
              className={clsx(
                "group relative w-full flex items-center gap-3 px-3 py-2 rounded-md text-sm transition-colors cursor-pointer",
                activeConfig === cfg
                  ? "bg-emerald-500/10 text-emerald-400 border border-emerald-500/20"
                  : "text-slate-400 hover:bg-white/5 hover:text-slate-200",
              )}
              onClick={() => setActiveConfig(cfg)}
            >
              <FileText className="w-4 h-4 shrink-0" />
              <span className="truncate flex-1">{cfg}</span>

              <button
                onClick={(e) => {
                  e.stopPropagation();
                  if (confirm(`确定要删除配置文件 ${cfg} 吗？`)) {
                    DeleteConfig(cfg)
                      .then(() => {
                        refreshConfigs();
                        if (activeConfig === cfg) {
                          setActiveConfig("");
                        }
                      })
                      .catch((err) => alert(err));
                  }
                }}
                className="opacity-0 group-hover:opacity-100 p-1 hover:bg-red-500/20 hover:text-red-400 text-slate-500 rounded transition-all"
                title="删除配置"
              >
                <Trash2 className="w-3.5 h-3.5" />
              </button>
            </div>
          ))}
        </div>
      </div>

      <div className="p-4 border-t border-white/5 mx-2">
        <div className="flex items-center gap-3 text-xs text-slate-500">
          <Database className="w-3 h-3" />
          v1.0.1
        </div>
      </div>
    </div>
  );
}
