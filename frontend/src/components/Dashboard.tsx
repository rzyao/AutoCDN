import { useState, useEffect, useRef } from "react";
import { Play, Terminal, Pause, CircleStop } from "lucide-react";
import { StartSpeedTest, StopSpeedTest } from "../../wailsjs/go/main/App";
import * as runtime from "../../wailsjs/runtime/runtime";
import clsx from "clsx";

declare global {
  interface Window {
    runtime: any;
  }
}

interface DashboardProps {
  activeConfig: string;
}

export default function Dashboard({ activeConfig }: DashboardProps) {
  const [running, setRunning] = useState(false);
  const [progress, setProgress] = useState(0);
  const [statusAction, setStatusAction] = useState("就绪");
  const [logs, setLogs] = useState<string[]>([]);
  const logsEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const cleanLog = runtime.EventsOn("log", (msg: string) => {
      setLogs((prev) => [...prev, `[日志] ${msg}`].slice(-100));
    });
    const cleanStatus = runtime.EventsOn("status", (msg: string) => {
      setStatusAction(msg);
      setLogs((prev) => [...prev, `[状态] ${msg}`].slice(-100));
    });
    const cleanProgress = runtime.EventsOn("progress", (data: any) => {
      if (data.total > 0) {
        setProgress((data.current / data.total) * 100);
      }
      setStatusAction(data.msg || "测速中...");
    });
    const cleanError = runtime.EventsOn("error", (msg: string) => {
      setLogs((prev) => [...prev, `[错误] ${msg}`].slice(-100));
    });

    return () => {
      cleanLog();
      cleanStatus();
      cleanProgress();
      cleanError();
    };
  }, []);

  useEffect(() => {
    logsEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [logs]);

  const handleStart = async (mode: "auto" | "manual") => {
    if (!activeConfig) return;
    setRunning(true);
    setLogs([]);
    setProgress(0);
    setStatusAction("正在初始化...");
    try {
      await StartSpeedTest(activeConfig, mode);
    } catch (e: any) {
      setLogs((prev) => [...prev, `[致命错误] ${e}`]);
    } finally {
      setRunning(false);
      setStatusAction("已完成 / 已停止");
      setProgress(100);
    }
  };

  const handleStop = async () => {
    try {
      await StopSpeedTest();
      // User logs will show "Stop command received"
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <div className="flex-1 p-8 flex flex-col h-full bg-slate-900 text-white overflow-hidden">
      <header className="mb-8 flex justify-between items-start">
        <div>
          <h2 className="text-3xl font-light text-white mb-2">开始优选</h2>
          <div className="text-slate-400">
            当前配置:{" "}
            <span className="text-emerald-400 font-mono">
              {activeConfig || "未选择"}
            </span>
          </div>
        </div>
      </header>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        {/* Auto Mode Card */}
        <div className="bg-gradient-to-br from-blue-600/20 to-blue-900/20 border border-blue-500/30 rounded-2xl p-6 relative overflow-hidden group hover:border-blue-400/50 transition-all">
          <div className="relative z-10">
            <h3 className="text-xl font-medium text-blue-100 mb-2">
              自动托管模式
            </h3>
            <p className="text-sm text-blue-200/60 mb-6">
              全自动执行延迟测速、优选 IP 并直接更新 Cloudflare DNS 记录。
            </p>
            <div className="flex gap-3">
              {!running ? (
                <button
                  onClick={() => handleStart("auto")}
                  disabled={!activeConfig}
                  className="bg-blue-500 hover:bg-blue-400 disabled:opacity-50 disabled:cursor-not-allowed text-white px-6 py-3 rounded-xl font-medium flex items-center gap-2 transition-all shadow-lg shadow-blue-500/20"
                >
                  <Play className="w-5 h-5 fill-current" />
                  一键优选
                </button>
              ) : (
                <button
                  onClick={handleStop}
                  className="bg-red-500 hover:bg-red-400 text-white px-6 py-3 rounded-xl font-medium flex items-center gap-2 transition-all shadow-lg shadow-red-500/20"
                >
                  <CircleStop className="w-5 h-5" />
                  停止任务
                </button>
              )}
            </div>
          </div>
          <div className="absolute -right-4 -bottom-4 opacity-10 group-hover:opacity-20 transition-opacity">
            <Play className="w-32 h-32" />
          </div>
        </div>

        {/* Manual Mode Card */}
        <div className="bg-slate-800/50 border border-white/10 rounded-2xl p-6 group hover:border-white/20 transition-all">
          <h3 className="text-xl font-medium text-slate-200 mb-2">专家模式</h3>
          <p className="text-sm text-slate-400 mb-6">
            仅执行测速并展示详细数据列表，需手动选择 IP 进行应用。
          </p>
          <div className="flex gap-3">
            {!running ? (
              <button
                onClick={() => handleStart("manual")}
                disabled={!activeConfig}
                className="bg-slate-700 hover:bg-slate-600 disabled:opacity-50 text-white px-6 py-3 rounded-xl font-medium flex items-center gap-2 transition-all"
              >
                <Terminal className="w-5 h-5" />
                仅测速分析
              </button>
            ) : (
              <button
                onClick={handleStop}
                className="bg-red-500/10 border border-red-500/50 text-red-400 hover:bg-red-500/20 px-6 py-3 rounded-xl font-medium flex items-center gap-2 transition-all"
              >
                <CircleStop className="w-5 h-5" />
                停止
              </button>
            )}
          </div>
        </div>
      </div>

      {/* Progress Section */}
      <div className="bg-slate-950 rounded-xl border border-white/5 flex-1 flex flex-col min-h-0">
        <div className="p-4 border-b border-white/5 flex justify-between items-center bg-white/5">
          <div className="flex items-center gap-2 text-sm font-mono text-slate-400">
            <Terminal className="w-4 h-4" />
            <span>控制台输出</span>
          </div>
          {running && (
            <div className="text-xs text-emerald-400 animate-pulse">
              {statusAction}
            </div>
          )}
        </div>

        {/* Progress Bar */}
        {running && (
          <div className="h-1 bg-slate-900 w-full">
            <div
              className="h-full bg-gradient-to-r from-blue-500 to-cyan-400 transition-all duration-300 ease-out"
              style={{ width: `${progress}%` }}
            />
          </div>
        )}

        <div className="flex-1 overflow-y-auto p-4 font-mono text-xs space-y-1">
          {logs.length === 0 && (
            <div className="text-slate-600 italic">准备就绪...</div>
          )}
          {logs.map((log, i) => (
            <div
              key={i}
              className="break-all whitespace-pre-wrap text-slate-300"
            >
              {log}
            </div>
          ))}
          <div ref={logsEndRef} />
        </div>
      </div>
    </div>
  );
}
