import { useState, useEffect } from "react";
import { Save, RefreshCw, AlertCircle } from "lucide-react";
import { LoadConfig, SaveConfig } from "../../wailsjs/go/main/App";
import { config as ConfigModels } from "../../wailsjs/go/models";
import clsx from "clsx";

interface ConfigEditorProps {
  activeConfig: string;
}

export default function ConfigEditor({ activeConfig }: ConfigEditorProps) {
  const [cfg, setCfg] = useState<ConfigModels.Config>(
    new ConfigModels.Config(),
  );
  const [loading, setLoading] = useState(false);
  const [status, setStatus] = useState("");

  useEffect(() => {
    if (activeConfig) {
      loadConfig();
    }
  }, [activeConfig]);

  const loadConfig = async () => {
    setLoading(true);
    setStatus("加载中...");
    try {
      const data = await LoadConfig(activeConfig);
      // Apply defaults to state if values are missing/zero
      if (!data.SpeedTest.Routines) data.SpeedTest.Routines = 200;
      if (!data.SpeedTest.PingTimes) data.SpeedTest.PingTimes = 4;
      if (!data.SpeedTest.TestCount) data.SpeedTest.TestCount = 10;
      if (!data.SpeedTest.DownloadTime) data.SpeedTest.DownloadTime = 10;
      if (!data.SpeedTest.TCPPort) data.SpeedTest.TCPPort = 443;
      if (!data.SpeedTest.MaxDelay) data.SpeedTest.MaxDelay = 200; // Default max delay
      if (!data.SpeedTest.MaxLossRate) data.SpeedTest.MaxLossRate = 0.2;
      // Note: MinDelay, MinSpeed can be 0, so checks are stricter if needed, but 0 is default anyway.

      setCfg(new ConfigModels.Config(data));
      setStatus("");
    } catch (e: any) {
      console.error(e);
      setStatus(`错误: ${e}`);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    setLoading(true);
    setStatus("保存中...");
    try {
      await SaveConfig(activeConfig, cfg);
      setStatus("已保存!");
      setTimeout(() => setStatus(""), 2000);
    } catch (e: any) {
      console.error(e);
      setStatus(`保存失败: ${e}`);
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (
    section: "Cloudflare" | "SpeedTest",
    field: string,
    value: any,
  ) => {
    const newCfg = new ConfigModels.Config(cfg);
    // @ts-ignore
    newCfg[section][field] = value;
    setCfg(newCfg);
  };

  const handleArrayChange = (
    section: "Cloudflare",
    field: string,
    value: string,
  ) => {
    const newCfg = new ConfigModels.Config(cfg);
    // @ts-ignore
    newCfg[section][field] = value.split("\n");
    setCfg(newCfg);
  };

  if (!activeConfig) {
    return (
      <div className="flex h-full items-center justify-center text-slate-500 gap-2">
        <AlertCircle className="w-5 h-5" />
        请在侧边栏选择一个配置文件
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full bg-slate-900 text-white p-8 overflow-hidden">
      <header className="mb-6 flex justify-between items-center">
        <div>
          <h2 className="text-3xl font-light text-white mb-2">配置编辑器</h2>
          <div className="text-slate-400">
            正在编辑:{" "}
            <span className="text-emerald-400 font-mono">{activeConfig}</span>
          </div>
        </div>
        <div className="flex items-center gap-3">
          <span className="text-sm text-emerald-400">{status}</span>
          <button
            onClick={loadConfig}
            disabled={loading}
            className="p-2 bg-slate-800 rounded-lg hover:bg-slate-700 transition"
            title="重新加载"
          >
            <RefreshCw className={clsx("w-5 h-5", loading && "animate-spin")} />
          </button>
          <button
            onClick={handleSave}
            disabled={loading}
            className="flex items-center gap-2 bg-blue-600 hover:bg-blue-500 text-white px-4 py-2 rounded-lg font-medium transition shadow-lg shadow-blue-600/20"
          >
            <Save className="w-4 h-4" />
            保存变更
          </button>
        </div>
      </header>

      <div className="flex-1 overflow-y-auto pr-4 space-y-8 pb-10">
        {/* Cloudflare Section */}
        <section>
          <h3 className="text-lg font-medium text-blue-400 mb-4 border-b border-blue-500/20 pb-2">
            Cloudflare API 设置
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-2">
              <label className="text-sm text-slate-400">API Key</label>
              <input
                type="password"
                value={cfg.Cloudflare?.APIKey || ""}
                onChange={(e) =>
                  handleChange("Cloudflare", "APIKey", e.target.value)
                }
                className="w-full bg-slate-950 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-blue-500 transition font-mono"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm text-slate-400">注册邮箱 (Email)</label>
              <input
                type="text"
                value={cfg.Cloudflare?.Email || ""}
                onChange={(e) =>
                  handleChange("Cloudflare", "Email", e.target.value)
                }
                className="w-full bg-slate-950 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-blue-500 transition"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm text-slate-400">
                区域 ID (Zone ID)
              </label>
              <input
                type="text"
                value={cfg.Cloudflare?.ZoneID || ""}
                onChange={(e) =>
                  handleChange("Cloudflare", "ZoneID", e.target.value)
                }
                className="w-full bg-slate-950 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-blue-500 transition font-mono"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm text-slate-400">
                主域名 (Zone Name)
              </label>
              <input
                type="text"
                value={cfg.Cloudflare?.ZoneName || ""}
                onChange={(e) =>
                  handleChange("Cloudflare", "ZoneName", e.target.value)
                }
                className="w-full bg-slate-950 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-blue-500 transition"
              />
            </div>
            <div className="space-y-2 row-span-2">
              <label className="text-sm text-slate-400">
                IPv4 域名列表 (每行一个)
              </label>
              <textarea
                value={cfg.Cloudflare?.Domains?.join("\n") || ""}
                onChange={(e) =>
                  handleArrayChange("Cloudflare", "Domains", e.target.value)
                }
                className="w-full h-32 bg-slate-950 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-blue-500 transition font-mono text-sm leading-relaxed"
              />
            </div>
            <div className="space-y-2 row-span-2">
              <label className="text-sm text-slate-400">
                IPv6 域名列表 (每行一个)
              </label>
              <textarea
                value={cfg.Cloudflare?.DomainIPv6s?.join("\n") || ""}
                onChange={(e) =>
                  handleArrayChange("Cloudflare", "DomainIPv6s", e.target.value)
                }
                className="w-full h-32 bg-slate-950 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-blue-500 transition font-mono text-sm leading-relaxed"
              />
            </div>
          </div>
        </section>

        {/* SpeedTest Section */}
        <section>
          <h3 className="text-lg font-medium text-emerald-400 mb-4 border-b border-emerald-500/20 pb-2">
            测速与筛选参数
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {/* 基础测速参数 */}
            <div className="space-y-2">
              <label className="text-sm text-slate-400">
                并发线程数 (Routines)
              </label>
              <input
                type="number"
                value={cfg.SpeedTest?.Routines}
                onChange={(e) =>
                  handleChange(
                    "SpeedTest",
                    "Routines",
                    parseInt(e.target.value),
                  )
                }
                className="w-full bg-slate-950 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-emerald-500 transition font-mono"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm text-slate-400">TCP 端口</label>
              <input
                type="number"
                value={cfg.SpeedTest?.TCPPort}
                onChange={(e) =>
                  handleChange("SpeedTest", "TCPPort", parseInt(e.target.value))
                }
                className="w-full bg-slate-950 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-emerald-500 transition font-mono"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm text-slate-400">Ping 次数</label>
              <input
                type="number"
                value={cfg.SpeedTest?.PingTimes}
                onChange={(e) =>
                  handleChange(
                    "SpeedTest",
                    "PingTimes",
                    parseInt(e.target.value),
                  )
                }
                className="w-full bg-slate-950 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-emerald-500 transition font-mono"
              />
            </div>

            {/* 筛选参数 */}
            <div className="space-y-2">
              <label className="text-sm text-slate-400">最大延迟 (ms)</label>
              <input
                type="number"
                value={cfg.SpeedTest?.MaxDelay}
                onChange={(e) =>
                  handleChange(
                    "SpeedTest",
                    "MaxDelay",
                    parseInt(e.target.value),
                  )
                }
                className="w-full bg-slate-950 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-emerald-500 transition font-mono"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm text-slate-400">最小延迟 (ms)</label>
              <input
                type="number"
                value={cfg.SpeedTest?.MinDelay || 0}
                onChange={(e) =>
                  handleChange(
                    "SpeedTest",
                    "MinDelay",
                    parseInt(e.target.value),
                  )
                }
                className="w-full bg-slate-950 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-emerald-500 transition font-mono"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm text-slate-400">
                最大丢包率 (0.0 - 1.0)
              </label>
              <input
                type="number"
                step="0.01"
                value={cfg.SpeedTest?.MaxLossRate || 0.2}
                onChange={(e) =>
                  handleChange(
                    "SpeedTest",
                    "MaxLossRate",
                    parseFloat(e.target.value),
                  )
                }
                className="w-full bg-slate-950 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-emerald-500 transition font-mono"
              />
            </div>

            {/* 下载测速参数 */}
            <div className="space-y-2">
              <label className="text-sm text-slate-400">下载测速目标数量</label>
              <input
                type="number"
                value={cfg.SpeedTest?.TestCount || 10}
                onChange={(e) =>
                  handleChange(
                    "SpeedTest",
                    "TestCount",
                    parseInt(e.target.value),
                  )
                }
                className="w-full bg-slate-950 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-emerald-500 transition font-mono"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm text-slate-400">
                单IP测速时长 (秒)
              </label>
              <input
                type="number"
                value={cfg.SpeedTest?.DownloadTime || 10}
                onChange={(e) =>
                  handleChange(
                    "SpeedTest",
                    "DownloadTime",
                    parseInt(e.target.value),
                  )
                }
                className="w-full bg-slate-950 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-emerald-500 transition font-mono"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm text-slate-400">
                最低下载速度 (MB/s)
              </label>
              <input
                type="number"
                step="0.1"
                value={cfg.SpeedTest?.MinSpeed || 0}
                onChange={(e) =>
                  handleChange(
                    "SpeedTest",
                    "MinSpeed",
                    parseFloat(e.target.value),
                  )
                }
                className="w-full bg-slate-950 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-emerald-500 transition font-mono"
              />
            </div>

            <div className="space-y-2 col-span-1 md:col-span-2 lg:col-span-3">
              <label className="text-sm text-slate-400">下载测速文件地址</label>
              <input
                type="text"
                value={cfg.SpeedTest?.SpeedTestURL || ""}
                onChange={(e) =>
                  handleChange("SpeedTest", "SpeedTestURL", e.target.value)
                }
                className="w-full bg-slate-950 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-emerald-500 transition font-mono"
              />
            </div>

            {/* 模式选择与文件配置 */}
            <div className="space-y-4 col-span-1 md:col-span-2 lg:col-span-3 border-t border-white/5 pt-6 mt-2">
              <h4 className="text-sm font-medium text-slate-500 uppercase tracking-wider mb-2">
                模式设置
              </h4>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="space-y-2">
                  <label className="text-sm text-slate-400">测速模式</label>
                  <select
                    value={cfg.SpeedTest?.TestType || "IPV4"}
                    onChange={(e) =>
                      handleChange("SpeedTest", "TestType", e.target.value)
                    }
                    className="w-full bg-slate-950 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-emerald-500 transition text-emerald-400 font-medium appearance-none cursor-pointer"
                  >
                    <option value="IPV4">IPv4 模式</option>
                    <option value="IPV6">IPv6 模式</option>
                  </select>
                </div>
                <div className="space-y-2">
                  <label className="text-sm text-slate-400">
                    IPv4 文件路径
                  </label>
                  <input
                    type="text"
                    value={cfg.SpeedTest?.IPv4File || "ip.txt"}
                    onChange={(e) =>
                      handleChange("SpeedTest", "IPv4File", e.target.value)
                    }
                    className={clsx(
                      "w-full bg-slate-950 border rounded-lg px-4 py-2 focus:outline-none transition font-mono",
                      cfg.SpeedTest?.TestType === "IPV6"
                        ? "border-white/5 text-slate-600"
                        : "border-white/10 focus:border-emerald-500",
                    )}
                    disabled={cfg.SpeedTest?.TestType === "IPV6"}
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm text-slate-400">
                    IPv6 文件路径
                  </label>
                  <input
                    type="text"
                    value={cfg.SpeedTest?.IPv6File || "ipv6.txt"}
                    onChange={(e) =>
                      handleChange("SpeedTest", "IPv6File", e.target.value)
                    }
                    className={clsx(
                      "w-full bg-slate-950 border rounded-lg px-4 py-2 focus:outline-none transition font-mono",
                      cfg.SpeedTest?.TestType !== "IPV6"
                        ? "border-white/5 text-slate-600"
                        : "border-white/10 focus:border-emerald-500",
                    )}
                    disabled={cfg.SpeedTest?.TestType !== "IPV6"}
                  />
                </div>
              </div>
            </div>

            <div className="flex flex-col gap-3 pt-6 col-span-1 md:col-span-2 lg:col-span-3">
              <h4 className="text-sm font-medium text-slate-500 uppercase tracking-wider">
                高级开关
              </h4>
              <div className="flex flex-wrap gap-6">
                <label className="flex items-center gap-3 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={cfg.SpeedTest?.Httping || false}
                    onChange={(e) =>
                      handleChange("SpeedTest", "Httping", e.target.checked)
                    }
                    className="w-5 h-5 accent-emerald-500 rounded bg-slate-700 border-none"
                  />
                  <span className="text-sm text-slate-300">
                    启用 HTTPing 模式
                  </span>
                </label>
                <label className="flex items-center gap-3 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={cfg.SpeedTest?.DisableDownload || false}
                    onChange={(e) =>
                      handleChange(
                        "SpeedTest",
                        "DisableDownload",
                        e.target.checked,
                      )
                    }
                    className="w-5 h-5 accent-emerald-500 rounded bg-slate-700 border-none"
                  />
                  <span className="text-sm text-slate-300">
                    禁用下载测速 (仅延迟)
                  </span>
                </label>
                <label className="flex items-center gap-3 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={cfg.SpeedTest?.TestAllIP || false}
                    onChange={(e) =>
                      handleChange("SpeedTest", "TestAllIP", e.target.checked)
                    }
                    className="w-5 h-5 accent-emerald-500 rounded bg-slate-700 border-none"
                  />
                  <span className="text-sm text-slate-300">
                    测速所有 IP (不随机)
                  </span>
                </label>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
}
