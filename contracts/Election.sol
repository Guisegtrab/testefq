// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Election {
    address public admin;

    modifier onlyAdmin() {
        require(msg.sender == admin, "Only admin can call this function");
        _;
    }

    modifier hasNotVoted() {
        require(!voters[msg.sender].voted, "You have already voted");
        _;
    }

    
    struct Candidate {
        uint id;
        string name;
        uint voteCount;
    }

    struct Voter {
        bool voted;
        uint voteIndex;
    }

    
    Candidate[] public candidates;

    // Mapeamento sobre os eleitores
    mapping(address => Voter) public voters;

    // Eventos
    event CandidateRegistered(uint indexed candidateId, string name);
    event Voted(address indexed voter, uint candidateId);
    event WinnerDeclared(uint winningCandidateId, string name, uint voteCount);

    // Registro de candidatos
    function registerCandidate(string memory _name) public onlyAdmin {
        
        uint candidateId = candidates.length;
        
        candidates.push(Candidate(candidateId, _name, 0));
        
        emit CandidateRegistered(candidateId, _name);
    }

    // Função para votação
    function vote(uint _candidateId) public hasNotVoted {
        require(_candidateId < candidates.length, "Invalid candidate ID");
        require(!voters[msg.sender].voted, "You have already voted");

        
        voters[msg.sender].voted = true;
        
        voters[msg.sender].voteIndex = _candidateId;
        
        candidates[_candidateId].voteCount++;

        
        emit Voted(msg.sender, _candidateId);
    }

    // Função para contar votos-vencedor
    function countVotes() public onlyAdmin {
        uint maxVoteCount = 0;
        uint winningCandidateId;

        
        for (uint i = 0; i < candidates.length; i++) {
            if (candidates[i].voteCount > maxVoteCount) {
                maxVoteCount = candidates[i].voteCount;
                winningCandidateId = i;
            }
        }

        // Emite o evento
        emit WinnerDeclared(winningCandidateId, candidates[winningCandidateId].name, maxVoteCount);
    }

    // Função de auditoria
    function audit() public view returns (uint[] memory, uint[] memory) {
        uint[] memory candidateIds = new uint[](candidates.length);
        uint[] memory voteCounts = new uint[](candidates.length);

        for (uint i = 0; i < candidates.length; i++) {
            candidateIds[i] = candidates[i].id;
            voteCounts[i] = candidates[i].voteCount;
        }

        return (candidateIds, voteCounts);
    }
}
