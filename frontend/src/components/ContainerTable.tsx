import React, { useState } from 'react';
import { Container, Ping } from '../api/containersApi.ts';
import ContainerRow from './ContainerRow.tsx';
import { Table } from 'react-bootstrap';

interface ContainerRowData {
  container: Container;
  ping: Ping | null;
}

interface ContainerTableProps {
  containers: Container[];
  onRefresh: () => void;
}

const ContainerTable: React.FC<ContainerTableProps> = ({ containers, onRefresh }) => {
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

  const toggleSortOrder = () => {
    setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
  };

  const flatRows: ContainerRowData[] = [];
  containers.forEach(container => {
    if (!container.pings || container.pings.length === 0) {
      flatRows.push({ container, ping: null });
    } else {
      container.pings.forEach(ping => {
        flatRows.push({ container, ping });
      });
    }
  });

  const sortedRows = flatRows.sort((a, b) => {
    const timeA = a.ping ? new Date(a.ping.last_ping_time).getTime() : 0;
    const timeB = b.ping ? new Date(b.ping.last_ping_time).getTime() : 0;
    return sortOrder === 'asc' ? timeA - timeB : timeB - timeA;
  });

  return (
    <div>
      <Table striped hover responsive>
        <thead>
          <tr>
            <th>ID</th>
            <th>IP</th>
            <th>Ping</th>
            <th
              style={{ cursor: 'pointer' }}
              onClick={toggleSortOrder}
              title="Сортировать по Last Ping Time"
            >
              Last Ping Time {sortOrder === 'asc' ? '↑' : '↓'}
            </th>
          </tr>
        </thead>
        <tbody>
          {sortedRows.map((row, index) => (
            <ContainerRow key={`${row.container.id}-${index}`} container={row.container} ping={row.ping} />
          ))}
        </tbody>
      </Table>
    </div>
  );
};

export default ContainerTable;
