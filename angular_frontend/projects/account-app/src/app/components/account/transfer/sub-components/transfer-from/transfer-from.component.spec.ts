import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TransferFromComponent } from './transfer-from.component';

describe('TransferFromComponent', () => {
  let component: TransferFromComponent;
  let fixture: ComponentFixture<TransferFromComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TransferFromComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TransferFromComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
