$(document).ready(function() {
    $("#getContactsBtn").click(function() {
        $.get("/api/v1/contacts", function(response) {
            if (response.status === 200) {
                displayContacts(response.data);
            } else {
                console.error('Error fetching contacts:', response.message);
            }
        }).fail(function(xhr, status, error) {
            console.error('Error fetching contacts:', error);
        });
    });

    function fetchContacts() {
        $.get("/api/v1/contacts", function(response) {
            if (response.status === 200) {
                displayContacts(response.data);
            } else {
                console.error('Error fetching contacts:', response.message);
            }
        }).fail(function(xhr, status, error) {
            console.error('Error fetching contacts:', error);
        });
    }

    function displayContacts(data) {
        var tableBody = $("#contactsTableBody");
        tableBody.empty();

        data.forEach(function(contact) {
            var row = $("<tr>");
            row.append($("<td>").text(contact.ID));
            row.append($("<td>").text(contact.Timestamp));
            row.append($("<td>").text(contact.Name));
            row.append($("<td>").text(contact.Email));
            row.append($("<td>").text(contact.Phone));
            row.append($("<td>").text(contact.ExternalID));
            tableBody.append(row);
        });
    }

    $("#addContactBtn").click(function() {
        $("#externalIdInput").val('');
    });

    $("#saveContactBtn").click(function() {
        var externalIdString = $("#externalIdInput").val();

        if (/^\d+$/.test(externalIdString) && !/\+/.test(externalIdString)) {
            var externalId = parseInt(externalIdString);

            $.post("/api/v1/contacts", JSON.stringify({ external_id: externalId }), function(response) {
                console.log('Contact added successfully:', response);
                $('#addContactModal').modal('hide');
                fetchContacts();
            }).fail(function(xhr, status, error) {
                console.error('Error adding contact:', error);
            });
        } else {
            alert('Invalid input. Please enter a non-negative whole number.');
        }
    });

    $(document).on('click', '#contactsTable tbody tr', function () {
        $('#contactsTable tbody tr').removeClass('selected');
        $(this).addClass('selected');
        var isSelected = $('#contactsTable tbody tr.selected').length > 0;
        $('#updateContactBtn').prop('disabled', !isSelected);
        $('#deleteContactBtn').prop('disabled', !isSelected);
    });

    $("#updateContactBtn").click(function() {
        var selectedRow = $('#contactsTable tbody tr.selected');
        var rowData = {
            ID: selectedRow.find('td:eq(0)').text(),
            Timestamp: selectedRow.find('td:eq(1)').text(),
            Name: selectedRow.find('td:eq(2)').text(),
            Email: selectedRow.find('td:eq(3)').text(),
            Phone: selectedRow.find('td:eq(4)').text(),
            ExternalID: selectedRow.find('td:eq(5)').text()
        };

        $('#updateContactModal').modal('show');
        $('#updateIdInput').val(rowData.ID);
        $('#updateTimestampInput').val(rowData.Timestamp);
        $('#updateNameInput').val(rowData.Name);
        $('#updateEmailInput').val(rowData.Email);
        $('#updatePhoneInput').val(rowData.Phone);
        $('#updateExternalIdInput').val(rowData.ExternalID);
    });

    $("#saveUpdatedContactBtn").click(function() {
        var updatedData = {
            // ID: $('#updateIdInput').val(),
            // Timestamp: $('#updateTimestampInput').val(),
            Name: $('#updateNameInput').val(),
            Email: $('#updateEmailInput').val(),
            Phone: $('#updatePhoneInput').val(),
            ExternalID: $('#updateExternalIdInput').val()
        };

        var contactId = $('#updateIdInput').val();
        var newContactId = parseInt(contactId);

        $.ajax({
            url: '/api/v1/contacts/' + newContactId,
            type: 'PUT',
            contentType: 'application/json',
            data: JSON.stringify(updatedData),
            success: function(response) {
                console.log('Contact updated successfully:', response);
                $('#updateContactModal').modal('hide');
                fetchContacts();
            },
            error: function(xhr, status, error) {
                console.error('Error updating contact:', error);
                alert('Failed to update contact. Please try again.');
            }
        });
    });

    $("#deleteContactBtn").click(function() {
        var isSelected = $('#contactsTable tbody tr.selected').length > 0;

        if (isSelected) {
            $('#confirmationModal').modal('show');
        } else {
            alert('Please select a contact to delete.');
        }
    });

    $("#confirmDeleteBtn").click(function() {
        var selectedRow = $('#contactsTable tbody tr.selected');
        var contactId = selectedRow.find('td:eq(0)').text();
        var newContactId = parseInt(contactId);

        $.ajax({
            url: '/api/v1/contacts/' + newContactId,
            type: 'DELETE',
            success: function(response) {
                console.log('Contact deleted successfully:', response);
                fetchContacts();
            },
            error: function(xhr, status, error) {
                console.error('Error deleting contact:', error);
                alert('Failed to delete contact. Please try again.');
            }
        });

        $('#confirmationModal').modal('hide');
    });
});
