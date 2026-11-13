import { useState, useEffect } from 'react';

const cache = {};

export function useStashData(url) {
  const [data, setData] = useState(cache[url] ?? null);
  const [loading, setLoading] = useState(!cache[url]);
  const [error, setError] = useState(false);
  const [prevUrl, setPrevUrl] = useState(url);

  // If the URL changes, sync state from the cache during render (the
  // React-recommended "adjusting state when a prop changes" pattern), so
  // the effect below only performs the async fetch.
  if (url !== prevUrl) {
    setPrevUrl(url);
    setData(cache[url] ?? null);
    setLoading(!cache[url]);
    setError(false);
  }

  useEffect(() => {
    if (cache[url]) {
      return;
    }

    fetch(url)
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then((json) => {
        cache[url] = json;
        setData(json);
        setLoading(false);
      })
      .catch((err) => {
        console.error(`Error fetching ${url}:`, err);
        if (!cache[url]) {
          setError(true);
          setLoading(false);
        }
      });
  }, [url]);

  return { data, loading, error };
}
