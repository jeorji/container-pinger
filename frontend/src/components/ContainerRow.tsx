import React, { useState } from 'react';
import { Container, Ping } from '../api/containersApi.ts';
import { Collapse } from 'react-bootstrap';

interface ContainerRowProps {
  container: Container;
  ping: Ping | null;
}

const ContainerRow: React.FC<ContainerRowProps> = ({ container, ping }) => {
  const [expanded, setExpanded] = useState(false);

  let rowClass = '';
  switch (container.state.toLowerCase()) {
    case 'created':
      rowClass = 'table-info';
      break;
    case 'running':
      rowClass = 'table-success';
      break;
    case 'paused':
      rowClass = 'table-warning';
      break;
    case 'exited':
      rowClass = 'table-secondary';
      break;
    case 'dead':
      rowClass = 'table-danger';
      break;
    default:
      rowClass = '';
  }

  const toggleExpand = () => setExpanded(!expanded);

  return (
    <>
      <tr className={rowClass} style={{ cursor: 'pointer' }} onClick={toggleExpand}>
        <td>{container.id}</td>
        <td>{ping ? ping.ip : 'N/A'}</td>
        <td>{ping ? `${ping.last_ping_latency} ms` : 'N/A'}</td>
        <td>{ping ? new Date(ping.last_ping_time).toLocaleString() : 'N/A'}</td>
      </tr>
      <tr>
        <td colSpan={4} className="p-0 border-0">
          <Collapse in={expanded}>
            <div className="p-3 bg-light">
              <h5>Информация о контейнере</h5>
              <ul className="list-unstyled mb-0">
                <li>
                  <strong>Name:</strong> {container.name}
                </li>
                <li>
                  <strong>Image:</strong> {container.image}
                </li>
                <li>
                  <strong>State:</strong> {container.state}
                </li>
                <li>
                  <strong>Status:</strong> {container.status}
                </li>
                <li>
                  <strong>Created at:</strong> {new Date(container.created_at).toLocaleString()}
                </li>
              </ul>
            </div>
          </Collapse>
        </td>
      </tr>
    </>
  );
};

export default ContainerRow;
