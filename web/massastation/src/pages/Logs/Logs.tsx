import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import axios from 'axios';
import { useQuery } from '@tanstack/react-query';
import { Button, Input, toast } from '@massalabs/react-ui-kit';
import { FiCopy, FiDownload, FiFileText, FiRefreshCw } from 'react-icons/fi';

import Intl from '@/i18n/i18n';
import { URL } from '@/const/url/url';
import { useResource } from '@/custom/api/useResource';

// Number of lines displayed in the viewer, which only shows the current session.
// The download always contains the whole log file, previous sessions included.
const DISPLAYED_LINES = 1000;

interface LogsParams {
  lines?: number;
  currentSession?: boolean;
}

async function fetchLogs(params: LogsParams = {}): Promise<string> {
  const { data } = await axios.get<string>(
    `${import.meta.env.VITE_BASE_API}/${URL.LOGS}`,
    {
      params,
      responseType: 'text',
      // keep the raw text, axios would otherwise try to parse it as JSON
      transformResponse: (raw) => raw,
    },
  );

  return data;
}

interface LogEntry {
  timestamp?: string;
  level?: string;
  message: string;
}

// Log lines are "<ISO8601 timestamp>\t<LEVEL>\t<message>".
// Other lines (e.g. stack traces) continue the previous entry.
function parseLine(line: string): LogEntry {
  const [timestamp, level, ...message] = line.split('\t');
  if (message.length && /^\d{4}-\d{2}-\d{2}T/.test(timestamp)) {
    return { timestamp, level, message: message.join('\t') };
  }
  return { message: line };
}

// "2026-10-29T12:41:44.041+0100" -> "12:41:44.041"
function shortTime(timestamp: string): string {
  return timestamp.slice(11, 23);
}

function levelClass(level?: string): string {
  switch (level) {
    case 'ERROR':
    case 'DPANIC':
    case 'PANIC':
    case 'FATAL':
      return 'text-s-error';
    case 'WARN':
      return 'text-s-warning';
    case 'DEBUG':
      return 'text-neutral';
    default:
      return 'text-s-info-1';
  }
}

function LogLine({ entry }: { entry: LogEntry }) {
  if (!entry.timestamp) {
    return (
      <div className="pl-[11.5rem] text-neutral whitespace-pre-wrap break-all">
        {entry.message}
      </div>
    );
  }

  return (
    <div className="flex gap-4 py-0.5 hover:bg-tertiary/30 rounded">
      <span className="w-24 shrink-0 text-neutral" title={entry.timestamp}>
        {shortTime(entry.timestamp)}
      </span>
      <span className={`w-14 shrink-0 font-semibold ${levelClass(entry.level)}`}>
        {entry.level}
      </span>
      <span
        className={`whitespace-pre-wrap break-all ${
          entry.level === 'DEBUG' ? 'text-neutral' : 'text-f-primary'
        }`}
      >
        {entry.message}
      </span>
    </div>
  );
}

export function Logs() {
  const [filter, setFilter] = useState('');
  const [isDownloading, setIsDownloading] = useState(false);
  const viewerRef = useRef<HTMLDivElement>(null);

  const { data: version } = useResource<string>(URL.VERSION);

  const {
    data: logs,
    isLoading,
    isError,
    isFetching,
    refetch,
  } = useQuery({
    queryKey: ['logs', DISPLAYED_LINES],
    queryFn: () =>
      fetchLogs({ lines: DISPLAYED_LINES, currentSession: true }),
    refetchOnWindowFocus: false,
  });

  const lines = useMemo(
    () => (logs ? logs.split('\n').filter((line) => line !== '') : []),
    [logs],
  );

  const filteredLines = useMemo(() => {
    const needle = filter.trim().toLowerCase();
    if (!needle) return lines;
    return lines.filter((line) => line.toLowerCase().includes(needle));
  }, [lines, filter]);

  // Most recent logs are at the bottom: keep them in view when logs are (re)loaded.
  useEffect(() => {
    if (viewerRef.current) {
      viewerRef.current.scrollTop = viewerRef.current.scrollHeight;
    }
  }, [logs]);

  const handleFilterChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement> | string) => {
      setFilter(typeof e === 'string' ? e : e.target.value);
    },
    [],
  );

  const handleCopy = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(filteredLines.join('\n'));
      toast.success(Intl.t('logs.copied'));
    } catch (error) {
      console.error(error);
      toast.error(Intl.t('logs.copy-error'));
    }
  }, [filteredLines]);

  const handleDownload = useCallback(async () => {
    setIsDownloading(true);
    try {
      const fullLogs = await fetchLogs();
      const blob = new Blob([fullLogs], { type: 'text/plain' });
      const timestamp = new Date().toISOString().replace(/[:.]/g, '-');

      const link = document.createElement('a');
      link.href = window.URL.createObjectURL(blob);
      link.download = `massastation-${version ?? 'unknown'}-${timestamp}.log`;
      link.click();
      window.URL.revokeObjectURL(link.href);
    } catch (error) {
      console.error(error);
      toast.error(Intl.t('logs.download-error'));
    } finally {
      setIsDownloading(false);
    }
  }, [version]);

  function renderViewerContent() {
    if (isLoading) {
      return <p className="text-neutral">{Intl.t('logs.loading')}</p>;
    }
    if (isError) {
      return <p className="text-s-error">{Intl.t('logs.load-error')}</p>;
    }
    if (!lines.length) {
      return <p className="text-neutral">{Intl.t('logs.empty')}</p>;
    }
    if (!filteredLines.length) {
      return <p className="text-neutral">{Intl.t('logs.no-match')}</p>;
    }
    return (
      <div>
        {filteredLines.map((line, index) => (
          <LogLine key={index} entry={parseLine(line)} />
        ))}
      </div>
    );
  }

  return (
    <div className="bg-primary text-f-primary pt-24 px-8 pb-12 min-h-screen">
      <div className="max-w-7xl mx-auto space-y-8">
        <div className="flex items-center gap-6 mb-8">
          <div className="p-4 bg-secondary rounded-xl">
            <FiFileText className="w-10 h-10 text-c-primary" />
          </div>
          <div className="flex flex-col justify-center">
            <h1 className="text-4xl font-bold text-f-primary cursor-default">
              {Intl.t('logs.title')}
            </h1>
            <p className="text-neutral mt-2 cursor-default">
              {Intl.t('logs.subtitle')}
            </p>
          </div>
        </div>

        <div className="bg-secondary/50 backdrop-blur-sm rounded-2xl border border-tertiary/20 p-8 space-y-6">
          <div className="flex flex-wrap items-center gap-4">
            <div className="flex-1 min-w-64">
              <Input
                name="logs-filter"
                value={filter}
                onChange={handleFilterChange}
                customClass="bg-tertiary/20 border border-tertiary/30 rounded-lg p-3 w-full
                  font-mono text-sm placeholder:text-neutral/60"
                placeholder={Intl.t('logs.search')}
              />
            </div>
            <Button
              onClick={() => refetch()}
              disabled={isFetching}
              className="flex items-center gap-2 px-6 py-3 bg-secondary hover:bg-c-hover text-f-primary
                         border border-tertiary/20 font-medium rounded-xl transition-all disabled:opacity-50"
            >
              <FiRefreshCw className={`w-5 h-5 ${isFetching ? 'animate-spin' : ''}`} />
              {Intl.t('logs.refresh')}
            </Button>
            <Button
              onClick={handleCopy}
              disabled={!filteredLines.length}
              className="flex items-center gap-2 px-6 py-3 bg-c-default hover:bg-c-hover
                         text-white font-medium rounded-xl transition-all disabled:opacity-50"
            >
              <FiCopy className="w-5 h-5" />
              {Intl.t('logs.copy')}
            </Button>
            <Button
              onClick={handleDownload}
              disabled={isDownloading}
              className="flex items-center gap-2 px-6 py-3 bg-c-default hover:bg-c-hover
                         text-white font-medium rounded-xl transition-all disabled:opacity-50"
            >
              <FiDownload className="w-5 h-5" />
              {Intl.t('logs.download')}
            </Button>
          </div>

          <p className="text-sm text-neutral cursor-default">
            {lines.length >= DISPLAYED_LINES
              ? Intl.t('logs.tail-info', { lines: String(DISPLAYED_LINES) })
              : Intl.t('logs.session-info')}
          </p>

          <div
            ref={viewerRef}
            data-testid="logs-viewer"
            className="bg-tertiary/20 border border-tertiary/30 rounded-xl p-4 h-[60vh]
                       overflow-auto font-mono text-xs leading-5"
          >
            {renderViewerContent()}
          </div>
        </div>
      </div>
    </div>
  );
}
