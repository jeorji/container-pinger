import React, { useEffect, useState } from 'react';
import ContainerTable from '../components/ContainerTable.tsx';
import { Container, fetchContainers } from '../api/containersApi.ts';
import { Spinner, Card, Form, InputGroup, Button } from 'react-bootstrap';

const Dashboard: React.FC = () => {
  const [containers, setContainers] = useState<Container[]>([]);
  const [error, setError] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);

  const [searchTerm, setSearchTerm] = useState<string>('');
  const [autoRefreshSeconds, setAutoRefreshSeconds] = useState<number>(0);
  const [countdown, setCountdown] = useState<number>(autoRefreshSeconds);

  const loadContainers = async () => {
    setLoading(true);
    setError('');
    try {
      const data = await fetchContainers();
      setContainers(data);
    } catch (err) {
      setError('Ошибка при загрузке данных. Пожалуйста, попробуйте снова.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadContainers();
  }, []);

  useEffect(() => {
    if (autoRefreshSeconds <= 0) return;
    setCountdown(autoRefreshSeconds);
    const interval = setInterval(() => {
      setCountdown(prev => {
        if (prev === 1) {
          loadContainers();
          return autoRefreshSeconds;
        }
        return prev - 1;
      });
    }, 1000);
    return () => clearInterval(interval);
  }, [autoRefreshSeconds]);

  const filteredContainers = containers.filter(container => {
    const term = searchTerm.toLowerCase();
    const matchBasic =
      container.id.toLowerCase().includes(term) ||
      container.state.toLowerCase().includes(term) ||
      container.name.toLowerCase().includes(term) ||
      container.image.toLowerCase().includes(term);
    let matchPing = false;
    if (container.pings && container.pings.length > 0) {
      matchPing = container.pings.some(ping =>
        ping.ip.toLowerCase().includes(term)
      );
    }
    return matchBasic || matchPing;
  });

  return (
    <div className="container mt-4">
      {error && (
        <div className="alert alert-danger" role="alert">
          {error}
        </div>
      )}

      <Card className="mb-4">
        <Card.Header>
          <div className="d-flex flex-wrap justify-content-between align-items-center">
            <div className="mb-2">
              <InputGroup>
                <InputGroup.Text>Поиск</InputGroup.Text>
                <Form.Control
                  type="text"
                  placeholder="Введите ID, имя или IP..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                />
              </InputGroup>
            </div>

            <div className="mb-2 d-flex align-items-center">
              <Form.Select
                style={{ width: '120px' }}
                value={autoRefreshSeconds}
                onChange={(e) => {
                  const val = Number(e.target.value);
                  setAutoRefreshSeconds(val);
                  setCountdown(val);
                }}
              >
                <option value={0}>Выкл</option>
                <option value={5}>5 сек</option>
                <option value={10}>10 сек</option>
                <option value={30}>30 сек</option>
                <option value={60}>1 мин</option>
              </Form.Select>
              {autoRefreshSeconds > 0 && (
                <div className="ms-3">
                  <strong>{countdown}</strong> сек.
                </div>
              )}
            </div>

            <div className="mb-2">
              <Button variant="primary" onClick={loadContainers} disabled={loading}>
                {loading ? (
                  <>
                    <Spinner animation="border" size="sm" /> Загрузка...
                  </>
                ) : (
                  'Обновить'
                )}
              </Button>
            </div>
          </div>
        </Card.Header>
      </Card>

      {loading && containers.length === 0 ? (
        <div className="text-center">
          <Spinner animation="border" role="status">
            <span className="visually-hidden">Загрузка...</span>
          </Spinner>
        </div>
      ) : filteredContainers.length === 0 ? (
        <div className="alert alert-info">
          Нет контейнеров, соответствующих запросу.
        </div>
      ) : (
        <ContainerTable containers={filteredContainers} onRefresh={loadContainers} />
      )}
    </div>
  );
};

export default Dashboard;

