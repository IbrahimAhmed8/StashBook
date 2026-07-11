import { useState, useEffect } from 'react';

const cache = {};

export function useStashData(url) {
  const [data, setData] = useState(cache[url] ?? null);
  const [loading, setLoading] = useState(!cache[url]);
  const [error, setError] = useState(false);

  useEffect(() => {
    if (cache[url]) {
      setData(cache[url]);
      setLoading(false);
      return;
    }

    setLoading(true);
    setError(false);

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
